#!/bin/bash
set -e
set -o pipefail

build_image="swagger-init"
build_dir="/go/src/git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/docs"

docker build -f docker/swagger-init.Dockerfile -t $build_image . && \
    docker run --rm -v `pwd`/docs:$build_dir $build_image && \
    chown -R $USER:$USER ./docs
