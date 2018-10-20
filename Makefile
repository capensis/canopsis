TAG:=develop

DISTRIBUTIONS=debian8,debian9,centos7
comma:=,

ifndef VERBOSE
.SILENT:
endif

docker_images:
	@docker build -f Dockerfile.core.template . -t canopsis/canopsis-core:${TAG}
	@docker build -f Dockerfile.prov . -t canopsis/canopsis-prov:${TAG}

packages: docker_images
	docker run -e FIX_OWNERSHIP=1000 -e CANOPSIS_PACKAGE_TAG=1.2 -e CANOPSIS_PACKAGE_REL=1 -v /home/clara/Projets/canopsis/build:/build -v /home/clara/Projets/canopsis/docker/packaging:/packages --entrypoint "/packages/package-debian-9.sh" --user=0 -it canopsis/canopsis-prov:develop

# Command to run inside the docker image
# cd /packages/ && FIX_OWNERSHIP=1000 CANOPSIS_PACKAGE_TAG=1.2 CANOPSIS_PACKAGE_REL=1 ./package-debian-9.sh

docker:
	@./debian9.sed Dockerfile.core.template | docker build -f - . -t canopsis

test:
	for distrib in $(subst ${comma}, ,${DISTRIBUTIONS})  ; do \
		sh $$distrib.sed Dockerfile.core.template ; \
	done
