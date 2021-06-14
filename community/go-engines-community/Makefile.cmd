include ../../Makefile.var
-include Makefile.var

BINARY:=${PROJECT}

build:
	@echo "Building ${PROJECT} project"
	@env GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=$(CGO) GO111MODULE=on go build ${GO_BUILD_FLAGS} ${GO_BUILD_CUSTOM_FLAGS} -o ${OUTPUT_DIR}/${BINARY} ${LDFLAGS}

clean:
	@echo "Cleaning ${PROJECT} project"
	@rm -rf ${OUTPUT_DIR}/${BINARY}

docker_image:
	@echo "Building ${PROJECT} project docker image"
	@export GOLANG_IMAGE_TAG=${GOLANG_IMAGE_TAG} && \
		export BINARY=${BINARY} && \
		export OUTPUT_DIR=${OUTPUT_DIR} && \
		export TAG=${TAG} && \
		cd ${ROOT_DIR} && ./docker/build-cmd-image.sh

test:
	@echo "Testing ${PROJECT} project"
	@go test ${GO_TEST_PARAM}
