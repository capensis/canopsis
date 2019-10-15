# go-engines-cat

This repository contains the sources for go engines and plugins that are part
of Canopsis CAT.

## Building

The plugins can be built with the command `make build`. This command assumes
that the `go-engines` repository is cloned in `../../canopsis/go-engines`. This
path can be changed in the file `go.mod`.

The plugins should be run and built with the same version of `go-engines`. A
plugin built with `go-engines` 3.20.0 will not work with an engine in version
3.20.1.
