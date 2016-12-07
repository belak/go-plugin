# go-plugin

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/belak/go-plugin)
[![Travis](https://img.shields.io/travis/belak/go-plugin.svg)](https://travis-ci.org/belak/go-plugin)
[![Coveralls](https://img.shields.io/coveralls/belak/go-plugin.svg)](https://coveralls.io/github/belak/go-plugin)

go-plugin is a wrapper around go-resolve which makes it easy to load
plugins that depend on each other based on a whitelist and blacklist.

# Types of Plugins

There are two types of plugins. Providers and Optional plugins.

Providers are always loaded by the system, regardless of the
whitelist/blacklist. They can still take arguments, even from optional
plugins.

Optional plugins are loaded based on the whitelist/blacklist.

# Plugin Requirements

Each plugin needs to be a function which runs the code to load that
plugin. Any type in the method signature will be handled using
injection. Any returned error will halt plugin loading and return that
error from Registry.Load.

# Important Dependencies

[go-resolve](https://github.com/belak/go-resolve) does most of the heavy lifting
by running a topological sort of all the plugins so they're loaded in the proper
order. It relies on [inject](https://github.com/codegangsta/inject) to do the
reflection work to actually load the plugins.

[glob](github.com/gobwas/glob) is used in the whitelist and blacklist to
determine which plugins should be loaded. There are a few tricks with the glob
syntax, so be sure to look here if you're having trouble.
