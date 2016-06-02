package plugin

import (
	"github.com/Unknwon/com"
	"github.com/gobwas/glob"
)

func matchingPlugins(all map[string]interface{}, rawWhitelist, rawBlacklist []string) ([]string, error) {
	var whitelist []glob.Glob
	var blacklist []glob.Glob

	// Compile all of the whitelist into globs
	for _, rawGlob := range rawWhitelist {
		g, err := glob.Compile(rawGlob, '.')
		if err != nil {
			return nil, err
		}
		whitelist = append(whitelist, g)
	}

	// Compile all of the blacklist into globs
	for _, rawGlob := range rawBlacklist {
		g, err := glob.Compile(rawGlob, '.')
		if err != nil {
			return nil, err
		}
		blacklist = append(blacklist, g)
	}

	// If the whitelist is empty, we want to match all plugins.
	if len(rawWhitelist) == 0 {
		whitelist = append(whitelist, glob.MustCompile("**", '.'))
	}

	var matching []string
	for item := range all {
		if matchesGloblist(item, whitelist) && !matchesGloblist(item, blacklist) {
			matching = com.AppendStr(matching, item)
		}
	}

	return matching, nil
}

// matchesGloblist is a simple function which tries an item against a
// slice of globs. It returns true if any of them match.
func matchesGloblist(item string, list []glob.Glob) bool {
	for _, glob := range list {
		if glob.Match(item) {
			return true
		}
	}

	return false
}
