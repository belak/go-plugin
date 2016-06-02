package plugin

import (
	"testing"

	"github.com/Unknwon/com"
	"github.com/stretchr/testify/assert"
)

func matchingPluginTestCase(t *testing.T, rawAll, whitelist, blacklist, expected []string) {
	// Transform from a convenient format into the one we use
	// internally to make it easier to test.
	all := make(map[string]interface{})
	for _, item := range rawAll {
		all[item] = true
	}

	// Grab all the matching plugins and ensure we don't get an error
	// and all the items in the slices are the same, even if they
	// aren't in the same order.
	res, err := matchingPlugins(all, whitelist, blacklist)
	assert.NoError(t, err)
	assert.True(t, com.CompareSliceStr(res, expected))
}

// This tests all the simple cases. There are many more, but they will
// be added when there are problems in the future.
var tests = []struct {
	All       []string
	Whitelist []string
	Blacklist []string
	Expected  []string
}{
	{
		All:       []string{"a.b.c"},
		Whitelist: []string{},
		Blacklist: []string{},
		Expected:  []string{"a.b.c"},
	},
	{
		All:       []string{"a.b.c", "b.c.d"},
		Whitelist: []string{"a.**"},
		Blacklist: []string{},
		Expected:  []string{"a.b.c"},
	},
	{
		All:       []string{"a.b.c", "b.c.d"},
		Whitelist: []string{},
		Blacklist: []string{"a.**"},
		Expected:  []string{"b.c.d"},
	},
	{
		All:       []string{"a.b.c", "a.c.d"},
		Whitelist: []string{"a.**"},
		Blacklist: []string{"a.b.**"},
		Expected:  []string{"a.c.d"},
	},
}

func TestMatchingPlugins(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		matchingPluginTestCase(
			t,
			test.All,
			test.Whitelist,
			test.Blacklist,
			test.Expected,
		)
	}

	_, err := matchingPlugins(
		map[string]interface{}{"a.b.c": true},
		[]string{"["},
		[]string{},
	)
	assert.Error(t, err)

	_, err = matchingPlugins(
		map[string]interface{}{"a.b.c": true},
		[]string{},
		[]string{"["},
	)
	assert.Error(t, err)
}
