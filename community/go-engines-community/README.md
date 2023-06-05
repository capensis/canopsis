# Canopsis Community Go engines

This directory contains the open-source “new generation” engines for Canopsis Community. They are written in Go and licensed under the [GNU AGPLv3](COPYING).

## Requirements

* A native build requires Requires [Go](https://golang.org/dl/) (see the `GOLANG_VERSION` variable in the [.env](../.env) for the exact version) and GNU Make.
* A Docker build requires Docker and GNU Make.

## Building

Run `make` to do a native build in your local environment. Resulting binaries will appear in the `build/` directory.

Run `make docker-images` to build the engines through Docker; the resulting images use the same format as our official containers.

Available targets can be listed with `make help`.

## Notes

* The Canopsis Community tree only builds the Canopsis Community parts. If you need a full set of Canopsis engines, you need to build Canopsis Community **AND** Canopsis Pro.
* Plugins should be run and built with the exact same version of Canopsis. A plugin built with `go-engines` 3.20.0 will not work with an engine in version 3.20.1.
