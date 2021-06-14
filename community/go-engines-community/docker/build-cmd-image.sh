#!/bin/bash
set -e
set -o pipefail

build_base="canopsis/go-build:${TAG}"
bin_base="canopsis/go-bin:${BINARY}-${TAG}"
bin_final="canopsis/${BINARY}:${TAG}"

echo "BUILD BASE IMAGE"

docker build \
    -f docker/Dockerfile.cmd.build \
    --build-arg GOLANG_IMAGE_TAG=${GOLANG_IMAGE_TAG} \
    -t ${build_base} \
    .

echo "BUILD BIN IMAGE"

[ "${BINARY}" = "canopsis-api" ] && sed -i -e 's/^#EXPOSE 8082/EXPOSE 8082/' docker/Dockerfile.cmd
docker build \
    -f docker/Dockerfile.cmd \
    --build-arg BASE=${build_base} \
    --build-arg BINARY_NAME=${BINARY} \
    --build-arg PROJECT_NAME=${BINARY} \
    --build-arg OUTPUT_DIR=${OUTPUT_DIR} \
    -t ${bin_base} \
    .
[ "${BINARY}" = "canopsis-api" ] && sed -i -e 's/^EXPOSE 8082/#EXPOSE 8082/' docker/Dockerfile.cmd

if [ -f cmd/${BINARY}/Dockerfile ]; then
    echo "BUILD FINAL IMAGE"

    docker build \
        -f cmd/${BINARY}/Dockerfile \
        --build-arg BASE=${bin_base} \
        -t ${bin_final} \
        .
else
    echo "TAG FINAL IMAGE"

    docker tag ${bin_base} ${bin_final}
fi

if [ "${BINARY}" = "canopsis-reconfigure" ];then
  echo "BUILD ${BINARY} CAT IMAGE"
  docker build \
      -f docker/Dockerfile.cmd-cat \
      --build-arg BASE=${bin_final} \
      -t "canopsis/${BINARY}-cat:${TAG}" \
      .
fi

