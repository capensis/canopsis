include Makefile.var

ARCHIVE_NAME:=canopsis-${TAG}-go-engines-community.tar.gz
ARCHIVE:=${RELEASES_DIR}/${ARCHIVE_NAME}
CHECKSUMS:=checksums.txt
SUBDIRS:=$(dir $(shell ls -d cmd/*/))
BUILD_DIRNAME:=canopsis-${TAG}
BUILD_OUTPUT_DIR?=${ROOT_DIR}/build
PARAM:=

ifndef VERBOSE
.SILENT:
endif

.PHONY: init test release clean full_clean build

.iter_project:
	@for dir in $(SUBDIRS) ; do \
		$(MAKE) -C $$dir $(MFLAGS) ${RULE} ${PARAM} -f ${ROOT_DIR}/Makefile.cmd; \
	done

release: clean
release: build
	@echo "Building $(ARCHIVE_NAME) for release"
	mkdir -p ${RELEASES_DIR}
	cd ${BUILD_OUTPUT_DIR}
	tar -czf $(ARCHIVE) -C ${BUILD_OUTPUT_DIR} --transform "s/^\./${BUILD_DIRNAME}/" .
	rm -rf ${BUILD_OUTPUT_DIR} && mkdir -p ${BUILD_OUTPUT_DIR}
	@echo "Build finished. Release file: $(ARCHIVE)"

clean: RULE = clean
clean: .iter_project

full_clean: RULE = clean
full_clean: .iter_project

init: 
	@echo "Init is deprecated and no longer needed to build"


test: 
	@echo "Running tests"
	go build -o ${BUILD_OUTPUT_DIR}/ready ./cmd/ready
	${BUILD_OUTPUT_DIR}/ready -retries 0 -timeout 2m
	cp ./cmd/canopsis-reconfigure/canopsis-core.toml.example ./canopsis.toml && go build -o ${BUILD_OUTPUT_DIR}/canopsis-reconfigure ./cmd/canopsis-reconfigure
	${BUILD_OUTPUT_DIR}/canopsis-reconfigure -conf ./canopsis.toml -mongoConf config/mongo
	# The options -p 1 -parallel 1 ensure that the tests are not run
	# concurrently, which is necessary for tests that use MongoDB, or
	# environment variables.
	go test -p 1 -parallel 1 ./...

build: RULE = build
build: PARAM:=OUTPUT_DIR=${BUILD_OUTPUT_DIR}/
build: .iter_project

docker_images: RULE = docker_image
docker_images: PARAM:=CLEAN_UP=false
docker_images: .iter_project

docker_test_image:
	docker build \
		-f docker/Dockerfile.test \
		--build-arg GOLANG_IMAGE_TAG=${GOLANG_IMAGE_TAG} \
		-t canopsis/go-test:latest \
		.

docker_release:
	@cd ${ROOT_DIR}
	docker build \
		-f docker/Dockerfile.release \
		--build-arg GOLANG_IMAGE_TAG=${GOLANG_IMAGE_TAG} \
		-t canopsis/go-release:${VERSION} \
		.
	docker run \
		-v ${ROOT_DIR}/releases/docker/:/releases/ \
		canopsis/go-release:${VERSION} \
		make release
	@echo "Release file: ${ROOT_DIR}/releases/docker/${ARCHIVE_NAME}"

help:
	@echo "Available targets: "
	@echo "   init: download go dependancies."
	@echo "   build: Builds the binary for production and writes the checksums of the generated file"
	@echo "   test: run test"
	@echo "   release: performs tests, build and creates an archive ready for deployment"
	@echo "   clean: remove binaries"
	@echo "   full_clean: call clean target and purge vendor directory and glide cache"
	@echo "   docker_images: builds docker image"
	@echo "   docker_release: builds docker image and make an archive ready to deploy"
	@echo "   docker_test_image: builds latest test docker image"
