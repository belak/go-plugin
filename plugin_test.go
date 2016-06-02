package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistrySimple(t *testing.T) {
	t.Parallel()

	r := NewRegistry()
	require.NotNil(t, r)

	// Ensure we can't register nil
	err := r.Register("hello.world", nil)
	assert.Error(t, err)
	assert.Empty(t, r.plugins)

	err = r.RegisterProvider(nil)
	assert.Error(t, err)
	assert.Empty(t, r.providers)

	// Now that we've made sure you can't register crappy values,
	// let's try to register the same plugin twice.
	err = r.Register("hello.world", func() {})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.plugins))

	err = r.Register("hello.world", func() {})
	assert.Error(t, err)
	assert.Equal(t, 1, len(r.plugins))

	// Ensure we can register something as a provider
	err = r.RegisterProvider(func() {})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.providers))
}

func TestRegistryLoad(t *testing.T) {
	var (
		needsInt    = func(int) {}
		providesInt = func() int { return 42 }
	)

	// Ensure the simple path works
	r := NewRegistry()
	require.NotNil(t, r)

	assert.NoError(t, r.RegisterProvider(providesInt))
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

	assert.NoError(t, r.RegisterProvider(providesInt))
	assert.NoError(t, r.RegisterProvider(providesInt))

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
