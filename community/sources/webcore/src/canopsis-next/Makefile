ifneq (3.81,$(firstword $(sort $(MAKE_VERSION) 3.81)))
$(error This makefile requires GNU Make >= 3.81)
endif

# Either community or pro
# note: the Canopsis frontend is identical between Community and Pro
CANOPSIS_EDITION = community

CURRENT_MAKEFILE := $(abspath $(firstword $(MAKEFILE_LIST)))
# Directory containing the main Makefile (the Makefile that was executed)
CURRENT_ROOT_DIR = $(realpath $(dir $(CURRENT_MAKEFILE)))
# Where to put the built binaries
CURRENT_BUILD_OUTPUT_DIR = $(CURRENT_ROOT_DIR)/build

# Path of the monorepo
MONOREPO_ROOT := $(realpath $(dir $(CURRENT_MAKEFILE))/../../../../../)
VPATH = $(dir $(CURRENT_MAKEFILE))

#
# BUILD FLAGS
#

SOURCE_DATE_EPOCH:=$(shell date -u +%s)
VERSION:=$(shell git -C "$(MONOREPO_ROOT)" for-each-ref refs/tags --sort=-authordate --format='%(refname:short)' --count=1 2>/dev/null || echo "0.0.0")
GITCOMMIT:=$(shell git -C "$(MONOREPO_ROOT)" log -n 1 --format=oneline --abbrev=12 --pretty="%h" 2>/dev/null || echo "unknown Git commit")

# Version of NodeJS to use
NODE_VERSION=14.18.3

# Yarn/NPM flags
export FORCE_COLOR=0
export NPM_CONFIG_COLOR=false
export NPM_CONFIG_PROGRESS=false
export NPM_CONFIG_SPIN=false
export NODE_OPTIONS="--max-old-space-size=6144" # canopsis/canopsis-pro#4066

#
# MAIN BUILD
#

.PHONY: all
all: $(CURRENT_BUILD_OUTPUT_DIR)/www

$(CURRENT_BUILD_OUTPUT_DIR)/www: jest.config.js package.json src/ vue.config.js
	if test "$$(node -v | sed 's/v//')" != "$(NODE_VERSION)" ;then\
		echo "ERROR: Couldn't find NodeJS $(NODE_VERSION) on your system";\
		exit 1;\
	fi
	yarn install
	yarn build \
		--env production \
		--dest "$@"

#
# INSTALLATION
#

.PHONY: install
install: $(CURRENT_BUILD_OUTPUT_DIR)/www
	@echo "==> Installing Canopsis frontend"
	@if test -e "$(DESTDIR)/opt/canopsis/srv/www" ; then\
		echo "ERROR: $(DESTDIR)/opt/canopsis/srv/www must not already exist";\
		exit 1;\
	fi
	install -d -m 755 -- "$(DESTDIR)/opt/canopsis/srv/www" "$(DESTDIR)/etc/nginx/conf.d"
	install -c -m 644 -p -- "$(MONOREPO_ROOT)/community/docker/files/etc/nginx/cors.inc" "$(DESTDIR)/etc/nginx/cors.inc"
	install -c -m 644 -p -- "$(MONOREPO_ROOT)/community/docker/files/etc/nginx/https.inc" "$(DESTDIR)/etc/nginx/https.inc"
	install -c -m 644 -p -- "$(MONOREPO_ROOT)/community/docker/files/etc/nginx/rundeck.inc" "$(DESTDIR)/etc/nginx/rundeck.inc"
	install -c -m 644 -p -- "$(MONOREPO_ROOT)/community/docker/files/etc/nginx/default.conf" "$(DESTDIR)/etc/nginx/conf.d/default.conf"
	cp -pR -- "$(CURRENT_BUILD_OUTPUT_DIR)/www" "$(DESTDIR)/opt/canopsis/srv/"

#
# DOCKER
#


#
# DOCKER FLAGS
#

export DOCKER_BUILDKIT:=1
export BUILDKIT_PROGRESS:=plain

# If set to 1, previous build caches will be ignored
DOCKER_NO_CACHE = 
DOCKER_TAG = $(VERSION)
DOCKERFILE_PATH = $(COMMUNITY_ROOT_DIR)/Dockerfile
DOCKER_IMAGE_BASE_PATH = docker.canopsis.net/docker/$(CANOPSIS_EDITION)
# The official nginx container is based on Debian unless stated otherwise in the tag.
DOCKER_NGINX_VERSION = 1.20.2

define DOCKER_BUILD_ARGS
--platform linux/amd64 \
--build-arg "NODE_VERSION=$(NODE_VERSION)-buster-slim" \
--build-arg "NGINX_VERSION=$(DOCKER_NGINX_VERSION)" 
endef

.PHONY: docker-check-version docker-images docker-nginx
docker-images: docker-nginx docker-nginx-openshift

# We require at least Docker 20.10 for the BuildKit features we use
docker-check-version:
	@if test "$$(docker version --format '{{.Server.Version}}' 2>/dev/null | awk -F'.' '{print 0+$$1 $$2}')" -lt 2010 ; then\
		echo "ERROR: Docker >= 20.10 is required";\
		exit 1;\
	fi

docker-nginx: docker-check-version
	docker build \
		-t "$(DOCKER_IMAGE_BASE_PATH)/nginx:$(DOCKER_TAG)" \
		$(if $(DOCKER_NO_CACHE),--no-cache,) \
		$(DOCKER_BUILD_ARGS) \
		--file "$(MONOREPO_ROOT)/community/docker/Dockerfile.nginx" \
		"$(MONOREPO_ROOT)"

docker-nginx-openshift: docker-check-version
	docker build \
		-t "$(DOCKER_IMAGE_BASE_PATH)/nginx-openshift:$(DOCKER_TAG)" \
		$(if $(DOCKER_NO_CACHE),--no-cache,) \
		$(DOCKER_BUILD_ARGS) \
		--file "$(MONOREPO_ROOT)/community/docker/Dockerfile.nginx" \
		--target nginx-openshift \
		"$(MONOREPO_ROOT)"

#
# HELP
#

.PHONY: help
help:
	@echo "Available targets are:"
	@$(MAKE) -pRrq -f "$(CURRENT_MAKEFILE)" : 2>/dev/null | awk -F':' '/^[a-zA-Z0-9][^$$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)print A[i]}' | sort -u
	@echo
	@echo "See also: $(CURRENT_ROOT_DIR)/README.md"

#
# CLEANUP
#

.PHONY: clean
clean:
	rm -rf -- "$(CURRENT_BUILD_OUTPUT_DIR)"
	rm -rf -- "./node_modules"
	rm -f -- "$(MONOREPO_ROOT)/.dockerignore"