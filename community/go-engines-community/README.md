# Canopsis Go-engines (Open Core)

This directory contains the open-source “new generation” engines, written in Go.

Licensed under the [GNU AGPLv3](COPYING).

## Requirements

Requires [Go](https://golang.org/dl/) and GNU Make.

See the `GOLANG_IMAGE_TAG` variable in [Makefile.var](Makefile.var) for the exact version.

## Building

Run `make` to build the binaries for your current environment. Resulting binaries will appear in the `build/` directory.

Run `make docker_images TAG="1.2.3" VERSION="1.2.3"` to build the engines through Docker images. Replace `1.2.3` with your current Git tag.

Run `make help` for more information.
