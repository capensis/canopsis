include ../../Makefile.var
-include Makefile.var

BINARY:=${PROJECT}

build:
	@echo "Building ${PROJECT} project"
	env CGO_ENABLED=1 GO111MODULE=on go build ${LDFLAGS} -trimpath -o "${OUTPUT_DIR}/${BINARY}"

clean:
	@echo "Cleaning ${PROJECT} project"
	rm -rf "${OUTPUT_DIR}/${BINARY}"

docker_image:
	@echo "Building ${PROJECT} project docker image"
	export GOLANG_IMAGE_TAG="${GOLANG_IMAGE_TAG}" BINARY="${BINARY}" TAG="${TAG}" && \
		cd "${ROOT_DIR}" && ./docker/build-cmd-image.sh

test:
	@echo "Testing ${PROJECT} project"
	go test ${GO_TEST_PARAM}
