TAG:=develop

DISTRIBUTIONS=debian8,debian9,centos7 # Every GNU/Linux distribution supported by Canopsis
# It's trick to allow subst to replace a comma.
comma:=,
DOCKER_DISTRIB="debian9" # The GNU/Linux distribution use as foundation for the official Canopsis Docker image

ifndef VERBOSE
.SILENT:
endif

docker_images:
	for distrib in $(subst ${comma}, ,${DISTRIBUTIONS})  ; do \
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
	for distrib in $(subst ${comma}, ,${DISTRIBUTIONS})  ; do \

		echo "*** Building " $$distrib " package"; \

		if [ "$$distrib" = ${DOCKER_DISTRIB} ]; then \
			export image_tag=${TAG}; \
		else \
			export image_tag=$$distrib-${TAG}; \
		fi; \

		docker run -e FIX_OWNERSHIP=``id -u`` -e CANOPSIS_PACKAGE_TAG=1.2 -e CANOPSIS_PACKAGE_REL=1 -v `pwd`/build:/build -v `pwd`/packaging:/packages --entrypoint "/packages/package-debian-9.sh" --user=0 -it canopsis/canopsis-prov:develop: \
	done


all: packages
