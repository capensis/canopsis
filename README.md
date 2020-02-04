# Canopsis Go-engines (Open Core)

This repository contains the open-source “new generation” engines, written in Go.

Licensed under the [GNU AGPLv3](COPYING).

## Requirements

Requires Go 1.12+, GNU Make, and Go Mod; see the `GOLANG_IMAGE_TAG` variable in [Makefile.var](Makefile.var) for the exact version.

## Building

Run `make build` to natively build the binaries in your current environment (Linux x86-64 only, for the moment). Resulting binaries will appear in the `build/` directory.

Run `make docker_images TAG="1.2.3" VERSION="1.2.3"` to build the engines through Docker images. Replace `1.2.3` with your current Git tag.

Run `make help` for more information.
