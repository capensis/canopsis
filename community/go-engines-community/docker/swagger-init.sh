#!/usr/bin/env bash
set -eu
set -o pipefail

build_image="swagger-init"
build_dir="/go/src/canopsis-community/docs"

docker build -f docker/swagger-init.Dockerfile -t "$build_image" . && \
    docker run --rm -v "$(pwd)/docs:$build_dir" "$build_image" && \
    chown -R "$(id -u):$(id -g)" ./docs
