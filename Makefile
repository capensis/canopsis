TAG:=develop

DISTRIBUTIONS=debian8,debian9,centos7 # Every GNU/Linux distribution supported by Canopsis
# It's trick to allow subst to replace a comma.
comma:=,
DOCKER_DISTRIB="debian9" # The GNU/Linux distribution use as foundation for the official Canopsis Docker image
PACKAGE_REV=""

ifndef VERBOSE
.SILENT:
endif

docker_images:
	for distrib in $(subst ${comma}, ,${DISTRIBUTIONS}) ; do \
		echo "*** Building " $$distrib; \
		if [ "$$distrib" = ${DOCKER_DISTRIB} ]; then \
			export image_tag=${TAG}; \
		else \
			export image_tag=$$distrib-${TAG}; \
		fi; \
		./$$distrib.sed Dockerfile.core.template | docker build -f - . -t canopsis/canopsis-core:$$image_tag; \
		sed 's|{{IMAGE_TAG}}|'$$image_tag'|' Dockerfile.prov.template ; \
		sed 's|{{IMAGE_TAG}}|'$$image_tag'|' Dockerfile.prov.template | docker build -f - . -t canopsis/canopsis-prov:$$image_tag; \
	done

packages: docker_images
	echo "Building packages" ; \
	for distrib in $(subst ${comma}, ,${DISTRIBUTIONS}) ; do \
		echo "*** Building " $$distrib " package"; \
		if [ "$$distrib" = ${DOCKER_DISTRIB} ]; then \
			export image_tag=${TAG}; \
		else \
			export image_tag=$$distrib-${TAG}; \
		fi; \
		docker run -e FIX_OWNERSHIP=``id -u`` -e CANOPSIS_PACKAGE_TAG=${TAG} -e CANOPSIS_PACKAGE_REL=${PACKAGE_REV} -v `pwd`/build:/build -v `pwd`/docker/packaging:/packages --entrypoint "/packages/package-"$$distrib".sh" --user=0 canopsis/canopsis-prov:develop ; \
	done
all: packages
