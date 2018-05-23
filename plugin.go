package plugin

import (
	"fmt"

	"github.com/belak/go-resolve"
	"github.com/codegangsta/inject"
)

// Registry represents a group of plugins which can be loaded. This
// is mostly a wrapper for go-resolve, so the same semantics apply
type Registry struct {
	plugins   map[string]interface{}
	providers map[string]interface{}
}

// NewRegistry creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins:   make(map[string]interface{}),
		providers: make(map[string]interface{}),
	}
}

// Copy creates a new instance of the registry which can be modified without
// changing the base. The general use case for this is instances where you want
// a global registry of plugins, but want to add separate providers for each
// instance of something, such as a bot.
func (r *Registry) Copy() *Registry {
	plugins := make(map[string]interface{})
	for k, v := range r.plugins {
		plugins[k] = v
	}

	providers := make(map[string]interface{})
	for k, v := range r.providers {
		providers[k] = v
	}

	return &Registry{
		plugins:   plugins,
		providers: providers,
	}
}

// Register simply ensures the provided factory is valid and adds the
// factory to a mapping under the given plugin name.
func (r *Registry) Register(name string, factory interface{}) error {
	if _, ok := r.plugins[name]; ok {
		return fmt.Errorf("A plugin with the name of %q is already loaded", name)
	}

	if _, ok := r.providers[name]; ok {
		return fmt.Errorf("A provider with the name of %q is already loaded", name)
	}

	err := resolve.EnsureValidFactory(factory)
	if err != nil {
		return err
	}

	r.plugins[name] = factory

	return nil
}

// RegisterProvider registers a factory as a part of the framework, so
// it will always be loaded.
func (r *Registry) RegisterProvider(name string, factory interface{}) error {
	if _, ok := r.providers[name]; ok {
		return fmt.Errorf("A provider with the name of %q is already loaded", name)
	}

	if _, ok := r.plugins[name]; ok {
		return fmt.Errorf("A plugin with the name of %q is already loaded", name)
	}

	err := resolve.EnsureValidFactory(factory)
	if err != nil {
		return err
	}

	r.providers[name] = factory

	return nil
}

// Load takes a whitelist and a blacklist in glob form and returns an
// inject.Injector representing the loaded plugins or an error
// representing why the loading failed.
func (r *Registry) Load(whitelist, blacklist []string) (inject.Injector, error) {
	res := resolve.NewResolver()

	// Add all non-plugin factories.
	for k, v := range r.providers {
		err := res.AddNode(k, v)
		if err != nil {
			return nil, err
		}
	}

	// Figure out which plugins match
	matching, err := matchingPlugins(r.plugins, whitelist, blacklist)
	if err != nil {
		return nil, err
	}

	// Now that we know all the plugins we want to load, we should do
	// that.
	for _, item := range matching {
		err := res.AddNode(item, r.plugins[item])
		if err != nil {
			return nil, err
		}
	}

	return res.Resolve()
}
