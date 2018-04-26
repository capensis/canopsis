#include Makefile.var

ifndef VERBOSE
.SILENT:
endif

#.PHONY: docker packages clean 
docker:
	@echo "build docker images"
	@export SYSBASE="debian-9"
	@./build-docker.sh ${TAG} ${START_BRANCH}
	@docker tag canopsis/canopsis-prov:${TAG} canopsis/canopsis-prov:latest
	@docker tag canopsis/canopsis-core:${TAG} canopsis/canopsis-core:latest


docker_push:
	@echo "build docker images and push them with tag and latest"
	$(MAKE) TAG=${TAG} START_BRANCH=${START_BRANCH} docker 
	@docker push 

docker_push: docker
docker_push: 
	@echo "push"

packages:
	@echo "build packages ${SYSBASE}"
	@./build-docker.sh ${TAG} ${START_BRANCH}
	@./build-packages.sh ${TAG}


all_packages:
	@echo "build all packages"
	$(MAKE) SYSBASE="centos-7" packages
	$(MAKE) SYSBASE="debian-8" packages
	$(MAKE) SYSBASE="debian-9" packages

packages_push:
	@echo "build and push all packages"
	$(MAKE) all_packages
	@echo "TODO push packages"

clean_docker:
	@echo "clean docker images"

clean_packages:
	@echo "clean packages"

clean:
	@echo "clean docker and packages"


help:
	@echo "Available targets: "
	@echo "docker: build docker images"
	@echo "docker_push: build docker images and push them"
	@echo "packages: build packages with sysbase"
	@echo "all_packages: build all packages"
	@echo "packages_push: build all packages and push them"
	@echo "clean: remove docker images and packages"
	@echo "clean_docker: remove docker images"
	@echo "clean_packages: remove packages"
	@echo "tag: tag canopsis project"
