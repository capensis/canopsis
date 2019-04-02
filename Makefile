ifndef VERBOSE # To use the verbose mode, add the VERBOSE=1 when using the Makefile
.SILENT:
endif

.DEFAULT_GOAL:=help

TAG:=develop

DISTRIBUTIONS=debian8,debian9,centos7 # Every GNU/Linux distribution supported by Canopsis
DOCKER_DISTRIB="debian9" # The GNU/Linux distribution use as foundation for the official Canopsis Docker image
PACKAGE_REV="1"
REGISTRY_NAME="git.canopsis.net"

# It's trick to allow subst to replace a comma.
.comma:=,


docker_images: DISTRIBUTIONS=debian9
docker_images:
	for distrib in $(subst ${.comma}, ,${DISTRIBUTIONS}) ; do \
                echo "*** Building " $$distrib; \
                if [ "$$distrib" = ${DOCKER_DISTRIB} ]; then \
                        export image_tag=${TAG}; \
                else \
                        export image_tag=$$distrib-${TAG}; \
                fi; \
                ./ci/$$distrib.sed ./ci/Dockerfile.template | docker build -f - . -t canopsis/canopsis:$$image_tag ; \
        done

docker_push: DISTRIBUTIONS=debian9
docker_push:
	for distrib in $(subst ${.comma}, ,${DISTRIBUTIONS}) ; do \
		echo "*** Push " $$distrib; \
		if [ "$$distrib" = ${DOCKER_DISTRIB} ]; then \
			export image_tag=${TAG}; \
		else \
			export image_tag=$$distrib-${TAG}; \
		fi; \
		docker tag canopsis/canopsis:$$image_tag ${REGISTRY_NAME}/canopsis/canopsis:$$image_tag ; \
		docker push ${REGISTRY_NAME}/canopsis/canopsis:$$image_tag ; \
	done

packages:
	echo "Building packages" ; \
        for distrib in $(subst ${.comma}, ,${DISTRIBUTIONS}) ; do \
                echo "*** Building " $$distrib " package"; \
                if [ "$$distrib" = ${DOCKER_DISTRIB} ]; then \
                        export image_tag=${TAG}; \
                else \
                        export image_tag=$$distrib-${TAG}; \
                fi; \
                docker run -e FIX_OWNERSHIP=`id -u`:`id -g` \
                           -e CANOPSIS_PACKAGE_TAG=${TAG} \
                           -e CANOPSIS_PACKAGE_REL=${PACKAGE_REV} \
                           -v `pwd`/build:/build \
                           -v `pwd`/docker/packaging:/packages \
                           --entrypoint "/packages/package-"$$distrib".sh" \
                           --user=0 canopsis/canopsis:$$image_tag ; \
        done

clean_images: DISTRIBUTIONS=debian9
clean_images:
	echo "Clean ...."; \
	for distrib in $(subst ${.comma}, ,${DISTRIBUTIONS}) ; do \
		if [ "$$distrib" = ${DOCKER_DISTRIB} ]; then \
			export image_tag=${TAG}; \
		else \
			export image_tag=$$distrib-${TAG}; \
		fi; \
		docker image rm canopsis/canopsis:$$image_tag ; \
	done	
	echo "Clean done."

all: packages

help:
	@echo "Available targets: "
	@echo "   * help: display this help."
	@echo "   * docker_images: Builds the docker image canopsis/canopsis."
	@echo "       - DISTRIBUTIONS: a coma separated list of GNU/Linux distribution."
	@echo "       Currently, debian8, debian9 and centos7 are supported. By default,"
	@echo "       only the debian9 images are build"
	@echo "       Example :"
	@echo "           - To build the debian9 and centos7 images"
	@echo "           make docker_images DISTRIBUTIONS=debian9,centos7"
	@echo "           - To build the debian8 images"
	@echo "           make docker_images DISTRIBUTIONS=debian8"
	@echo "   * clean: Clean docker images."
	@echo "       - DISTRIBUTIONS: a coma separated list of GNU/Linux distribution."
	@echo "       Currently, debian8, debian9 and centos7 are supported. By default,"
	@echo "       only the debian9 images are clean"
	@echo "       Example:"
	@echo "           - To clean the debian9 and centos7 images"
	@echo "           make docker_images DISTRIBUTIONS=debian9,centos7"
	@echo "           - To clean the debian8 images"
	@echo "           make docker_images DISTRIBUTIONS=debian8"
	@echo "   * packages: Builds the canopsis-core package. The package will be stored"
	@echo "   in the 'build' directory."
	@echo "       - DISTRIBUTIONS: a coma separated list of GNU/Linux distribution."
	@echo "       Currently, debian8, debian9 and centos7 are supported."
	@echo "       Example :"
	@echo "           - To build the debian9 and centos7 images"
	@echo "           make packages DISTRIBUTIONS=debian9,centos7"
	@echo "           - To build the debian8 images"
	@echo "           make packages DISTRIBUTIONS=debian8"
