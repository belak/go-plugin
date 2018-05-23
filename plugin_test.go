package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	needsInt    = func(int) {}
	providesInt = func() int { return 42 }
)

func TestRegistrySimple(t *testing.T) {
	t.Parallel()

	r := NewRegistry()
	require.NotNil(t, r)

	// Ensure we can't register nil
	err := r.Register("hello.world", nil)
	assert.Error(t, err)
	assert.Empty(t, r.plugins)

	err = r.RegisterProvider("core/nil", nil)
	assert.Error(t, err)
	assert.Empty(t, r.providers)

	// Ensure we can register something as a provider
	err = r.RegisterProvider("core/empty", func() {})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(r.plugins))
	assert.Equal(t, 1, len(r.providers))

	// Ensure we can't register a provider over another
	err = r.RegisterProvider("core/empty", func() {})
	assert.Error(t, err)
	assert.Equal(t, 0, len(r.plugins))
	assert.Equal(t, 1, len(r.providers))

	// Now that we've made sure you can't register crappy values,
	// let's try to register the same plugin twice.
	err = r.Register("hello.world", func() {})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.plugins))
	assert.Equal(t, 1, len(r.providers))

	err = r.Register("hello.world", func() {})
	assert.Error(t, err)
	assert.Equal(t, 1, len(r.plugins))
	assert.Equal(t, 1, len(r.providers))

	// We also shouldn't be able to register a provider over a plugin.
	err = r.RegisterProvider("hello.world", func() {})
	assert.Error(t, err)
	assert.Equal(t, 1, len(r.plugins))
	assert.Equal(t, 1, len(r.providers))

	// We also shouldn't be able to register a plugin over a provider.
	err = r.Register("core/empty", func() {})
	assert.Error(t, err)
	assert.Equal(t, 1, len(r.plugins))
	assert.Equal(t, 1, len(r.providers))
}

func TestRegistryLoad(t *testing.T) {
	t.Parallel()

	// Ensure the simple path works
	r := NewRegistry()
	require.NotNil(t, r)

	assert.NoError(t, r.RegisterProvider("core/int", providesInt))
	assert.NoError(t, r.Register("requires.int", needsInt))

	_, err := r.Load(nil, nil)
	assert.NoError(t, err)

	// Ensure loading with an invalid glob errors
	_, err = r.Load([]string{"["}, nil)
	assert.Error(t, err)

	// Ensure we can't load the same provider multiple times (because
	// of overlapping return types)
	r = NewRegistry()
	require.NotNil(t, r)

	assert.NoError(t, r.RegisterProvider("core/int1", providesInt))
	assert.NoError(t, r.RegisterProvider("core/int2", providesInt))

	_, err = r.Load(nil, nil)
	assert.Error(t, err)

	// Similar to the last check, we want to ensure that adding
	// multiple plugins which return the same values will error.
	r = NewRegistry()
	require.NotNil(t, r)

	assert.NoError(t, r.Register("provides.int.1", providesInt))
	assert.NoError(t, r.Register("provides.int.2", providesInt))

	_, err = r.Load(nil, nil)
	assert.Error(t, err)
}

func TestRegistryCopy(t *testing.T) {
	t.Parallel()

	r := NewRegistry()
	assert.NoError(t, r.Register("requires.int", needsInt))
	assert.NoError(t, r.RegisterProvider("core/int", providesInt))

	rcopy := r.Copy()
	assert.NoError(t, rcopy.Register("requires.int.2", needsInt))
	assert.NoError(t, rcopy.RegisterProvider("core/int2", providesInt))

	_, err := r.Load(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.plugins))
	assert.Equal(t, 1, len(r.providers))

	_, err = rcopy.Load(nil, nil)
	assert.Error(t, err)
	assert.Equal(t, 2, len(rcopy.plugins))
	assert.Equal(t, 2, len(rcopy.providers))
}
