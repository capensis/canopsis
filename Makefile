TAG:=develop

DISTRIB=debian8,debian9,centos7

ifndef VERBOSE
.SILENT:
endif

build_images:
	@docker build -f Dockerfile.core.template . -t canopsis/canopsis-core:${TAG}
	@docker build -f Dockerfile.prov . -t canopsis/canopsis-prov:${TAG}

build_packages: build_images
	docker run -v /home/clara/Projets/canopsis/docker/packaging:/packages --entrypoint "/bin/bash" --user=0 -it canopsis/canopsis-core:develop
