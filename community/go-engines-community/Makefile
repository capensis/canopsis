MAIN_MAKEFILE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURRENT_MAKEFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
# Path of the monorepo
MONOREPO_ROOT := $(shell realpath --relative-to=. $(dir $(CURRENT_MAKEFILE))/../..)
# Directory containing the main Makefile (the Makefile that was executed)
ROOT_DIR = $(shell realpath --relative-to=. $(dir $(MAIN_MAKEFILE)))
VPATH = $(dir $(MAIN_MAKEFILE))
# Where to put the built binaries
OUTPUT_DIR = $(ROOT_DIR)/build

CANOPSIS_EDITION = community

# List of cmds/engines to build and install (eg : canopsis-api connector-junit ...)
CMDS = $(notdir $(wildcard $(ROOT_DIR)/cmd/*))
# List of cmds binaries paths (eg : build/cmd/canopsis-api build/cmd/connector-junit ...)
CMDS_BINARY_TARGETS = $(addprefix $(OUTPUT_DIR)/cmd/,$(CMDS))

# List of plugins to build and install (eg : datasource/mongo.so)
PLUGINS = $(addsuffix .so,$(patsubst $(ROOT_DIR)/plugins/%,%,$(wildcard $(ROOT_DIR)/plugins/*/*)))
PLUGINS_BINARY_TARGETS = $(addprefix $(OUTPUT_DIR)/plugins/,$(PLUGINS))

# List of docker containers targets
DOCKER_TARGETS = $(addprefix docker_,$(CMDS))

# BUILD FLAGS
BUILDDATE:=$(shell date "+%Y-%m-%dT%H:%M:%S%z")
LINKPKG="git.canopsis.net/canopsis/canopsis-$(CANOPSIS_EDITION)/$(CANOPSIS_EDITION)/go-engines-$(CANOPSIS_EDITION)"
VERSION:=$(shell git -C $(MONOREPO_ROOT) for-each-ref refs/tags --sort=-authordate --format='%(refname:short)' --count=1)
GITCOMMIT:=$(shell git -C $(MONOREPO_ROOT) rev-parse HEAD)
GITBRANCH:=$(shell git -C $(MONOREPO_ROOT) rev-parse --abbrev-ref HEAD)
#Metadata for the binaries
CUSTOM_LDFLAGS=
LDFLAGS:=-ldflags "\
	-X ${LINKPKG}.BuildGitCommit=${GITCOMMIT} \
	-X ${LINKPKG}.BuildGitBranch=${GITBRANCH} \
	-X ${LINKPKG}.BuildDate=${BUILDDATE} ${CUSTOM_LDFLAGS}"

# Use the trimpath option so that the logs do not contain the absolute paths to
# the source files.
# This makes the log lines shorter, and independent from the build machine.
GO_BUILD_FLAGS:= -mod=mod #-gcflags="all=-trimpath \"${GOPATH}/src\""
GO_BUILD_CUSTOM_FLAGS=

# Export configuration env variables to the go compiler
export GOOS:=linux
export GOARCH:=amd64
export CGO_ENABLED:=1
export GO111MODULE:=on

#Docker configuration
export DOCKER_BUILDKIT = 1

#Docker variables
DOCKER_NO_CACHE = # If set to 1 previous build caches will be ignored
DOCKER_TAG = $(VERSION)
DOCKER_ALPINE_VERSION = 3.13
DOCKER_GOLANG_VERSION = 1.16.4
DOCKERFILE_PATH = $(dir $(CURRENT_MAKEFILE))Dockerfile
DOCKER_IMAGE_BASE_PATH = docker.canopsis.net/docker/$(CANOPSIS_EDITION)

.PHONY: build docker_images
build: $(CMDS_BINARY_TARGETS) $(PLUGINS_BINARY_TARGETS)
docker_images: $(DOCKER_TARGETS)

#
# BINARY TARGETS
#

# Make make treat cmds as viable build targets (for single-engine build)
$(CMDS_BINARY_TARGETS):
$(PLUGINS_BINARY_TARGETS):
# Tell make where to find the binaries


# Plugin build job template
$(OUTPUT_DIR)/plugins/%.so: $(ROOT_DIR)/plugins/% $(ROOT_DIR)/plugins/%/*
	@echo "Building $@"
	@cd $(ROOT_DIR);\
	go build $(GO_BUILD_FLAGS) $(GO_BUILD_CUSTOM_FLAGS) $(LDFLAGS) \
		-buildmode=plugin \
		-o "$(abspath $@)" "$(abspath $<)"


# Build job template
$(OUTPUT_DIR)/cmd/%: $(ROOT_DIR)/cmd/% $(ROOT_DIR)/cmd/%/*
	@echo "Building $@"
	@cd $(ROOT_DIR);\
	go build $(GO_BUILD_FLAGS) $(GO_BUILD_CUSTOM_FLAGS) $(LDFLAGS) \
		-o "$(abspath $@)" "$(abspath $<)"

# To add a custom project-specific build flag
#engine-action: GO_BUILD_CUSTOM_FLAGS=--your-custom-flag

#
# DOCKER IMAGES TARGETS
#

$(MONOREPO_ROOT)/.dockerignore:
	# I love Docker https://github.com/moby/moby/issues/20944
	cd $(MONOREPO_ROOT);\
	echo > $(MONOREPO_ROOT)/.dockerignore;\
	set -x;\
	for file in $$(find -name .dockerignore);\
	do \
		sed -E 's;^(!?)(.*)$$;\1'$$(dirname "$$file")'/\2;' "$$file" >> .dockerignore;\
	done

canopsis.toml:
	mkdir -p $(OUTPUT_DIR)
	cp $(dir $(CURRENT_MAKEFILE))/config/canopsis.${CANOPSIS_EDITION}.toml $(OUTPUT_DIR)/canopsis.toml

.PHONY: docker_alpine_updated docker_builder $(DOCKER_TARGETS)
docker_alpine_updated: $(MONOREPO_ROOT)/.dockerignore
	docker build \
		$(if $(DOCKER_NO_CACHE),--no-cache,--cache-from $(DOCKER_IMAGE_BASE_PATH)/alpine_updated:latest) \
		-t $(DOCKER_IMAGE_BASE_PATH)/alpine_updated:latest \
		--target alpine_updated \
		--build-arg "ALPINE_VERSION=${DOCKER_ALPINE_VERSION}" \
		--build-arg "CANOPSIS_EDITION=${CANOPSIS_EDITION}" \
		--build-arg "GOLANG_VERSION=${DOCKER_GOLANG_VERSION}" \
		--file "$(DOCKERFILE_PATH)" \
		"$(MONOREPO_ROOT)"

docker_builder: $(MONOREPO_ROOT)/.dockerignore docker_alpine_updated
	docker build \
		-t $(DOCKER_IMAGE_BASE_PATH)/canopsis_go_bin:latest \
		--target builder \
		--cache-from $(DOCKER_IMAGE_BASE_PATH)/alpine_updated:latest \
		$(if $(DOCKER_NO_CACHE),,--cache-from $(DOCKER_IMAGE_BASE_PATH)/canopsis_go_bin:latest) \
		--build-arg "ALPINE_VERSION=${DOCKER_ALPINE_VERSION}" \
		--build-arg "CANOPSIS_EDITION=${CANOPSIS_EDITION}" \
		--build-arg "GOLANG_VERSION=${DOCKER_GOLANG_VERSION}" \
		--file "$(DOCKERFILE_PATH)" \
		"$(MONOREPO_ROOT)"

$(DOCKER_TARGETS): TARGET=final_container
$(DOCKER_TARGETS): $(MONOREPO_ROOT)/.dockerignore docker_alpine_updated docker_builder
	docker build -t "$(DOCKER_IMAGE_BASE_PATH)/$(subst docker_,,$@):$(DOCKER_TAG)" \
		--build-arg CMD="$(subst docker_,,$@)" \
		--target $(TARGET) \
		--cache-from $(DOCKER_IMAGE_BASE_PATH)/alpine_updated:latest \
		--cache-from $(DOCKER_IMAGE_BASE_PATH)/canopsis_go_bin:latest \
		$(if $(DOCKER_NO_CACHE),,--cache-from "$(subst docker_,,$@):$(DOCKER_TAG)") \
		--build-arg "ALPINE_VERSION=${DOCKER_ALPINE_VERSION}" \
		--build-arg "CANOPSIS_EDITION=${CANOPSIS_EDITION}" \
		--build-arg "GOLANG_VERSION=${DOCKER_GOLANG_VERSION}" \
		--file "$(DOCKERFILE_PATH)" \
		"$(MONOREPO_ROOT)"

docker_canopsis-api: TARGET=canopsis-api

#
# Install targets
#

.PHONY: install

# You can modify the content of the CMDS and PLUGINS variables to change what will be built/installed
# Example : make install DESTDIR=/tmp/installroot PLUGINS=datasource/mongo.so CMDS="canopsis-api connector-junit"
# will only install the mongo datasource plugin, canopsis-api and connector-junit
install: canopsis.toml $(CMDS_BINARY_TARGETS) $(PLUGINS_BINARY_TARGETS)
	# Install cmds
	install -d $(DESTDIR)$(PREFIX)/opt/canopsis/bin
	$(foreach cmd,$(CMDS),install -m 755 $(OUTPUT_DIR)/cmd/$(cmd) $(DESTDIR)$(PREFIX)/opt/canopsis/bin/;)
	# Install plugins
	$(foreach plugin,$(PLUGINS),install -D -m 755 $(OUTPUT_DIR)/plugins/$(plugin) $(DESTDIR)$(PREFIX)/opt/canopsis/lib/go/plugins/$(plugin);)
	# Install configuration file
	install -d $(DESTDIR)$(PREFIX)/opt/canopsis/etc
	install -m 644 $(OUTPUT_DIR)/canopsis.toml $(DESTDIR)$(PREFIX)/opt/canopsis/etc/

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
	rm -f $(MONOREPO_ROOT)/.dockerignore
