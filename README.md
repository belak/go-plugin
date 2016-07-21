# go-plugin

[![Build Status](https://travis-ci.org/belak/go-plugin.svg?branch=master)](https://travis-ci.org/belak/go-plugin)
[![Coverage Status](https://coveralls.io/repos/github/belak/go-plugin/badge.svg?branch=master)](https://coveralls.io/github/belak/go-plugin?branch=master)

go-plugin is a wrapper around go-resolve which makes it easy to load plugins
based on a whitelist and blacklist which depend on each other.

# Important Dependencies

[go-resolve](https://github.com/belak/go-resolve) does most of the heavy lifting
by running a topological sort of all the plugins so they're loaded in the proper
order. It relies on [inject](https://github.com/codegangsta/inject) to do the
reflection work to actually load the plugins.

[glob](github.com/gobwas/glob) is used in the whitelist and blacklist to
determine which plugins should be loaded. There are a few tricks with the glob
syntax, so be sure to look here if you're having trouble.
