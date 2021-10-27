MAIN_MAKEFILE := $(abspath $(firstword $(MAKEFILE_LIST)))
CURRENT_MAKEFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
# Path of the monorepo
MONOREPO_ROOT := $(realpath $(dir $(CURRENT_MAKEFILE))/../..)
# Directory containing the main Makefile (the Makefile that was executed)
ROOT_DIR = $(realpath $(dir $(MAIN_MAKEFILE)))
VPATH = $(dir $(MAIN_MAKEFILE))
# Where to put the built binaries
BUILD_DIR = $(ROOT_DIR)/build

# Either community or pro
CANOPSIS_EDITION = community

# List of cmds/engines to build and install (eg: canopsis-api connector-junit...)
CMDS = $(notdir $(wildcard $(ROOT_DIR)/cmd/*))
# List of cmd binary paths (eg: build/cmd/canopsis-api build/cmd/connector-junit...)
CMDS_BINARY_TARGETS = $(addprefix $(BUILD_DIR)/cmd/,$(CMDS))

# List of plugins to build and install (eg: datasource/mongo.so)
PLUGINS = $(addsuffix .so,$(patsubst $(ROOT_DIR)/plugins/%,%,$(wildcard $(ROOT_DIR)/plugins/*/*)))
PLUGINS_BINARY_TARGETS = $(addprefix $(BUILD_DIR)/plugins/,$(PLUGINS))

# List of Docker containers targets
DOCKER_TARGETS = $(addprefix docker-,$(CMDS))

#
# BUILD FLAGS
#

BUILDDATE:=$(shell date -u "+%Y-%m-%dT%H:%M:%S%z")
LINKPKG=git.canopsis.net/canopsis/canopsis-$(CANOPSIS_EDITION)/$(CANOPSIS_EDITION)/go-engines-$(CANOPSIS_EDITION)/lib/canopsis
VERSION:=$(shell git -C "$(MONOREPO_ROOT)" for-each-ref refs/tags --sort=-authordate --format='%(refname:short)' --count=1 2>/dev/null || echo "unknown Git version")
GITCOMMIT:=$(shell git -C "$(MONOREPO_ROOT)" log -n 1 --format=oneline --abbrev-commit --pretty="%h" 2>/dev/null || echo "unknown Git commit")
# Metadata for the binaries
CUSTOM_LDFLAGS=
LDFLAGS:=-ldflags "\
	-X '${LINKPKG}.BuildGitCommit=${GITCOMMIT}' \
	-X '${LINKPKG}.BuildDate=${BUILDDATE}' ${CUSTOM_LDFLAGS}"

# -trimpath makes the log lines shorter and independent from the current machine.
GO_BUILD_FLAGS:= -trimpath
GO_BUILD_CUSTOM_FLAGS=

# Export configuration env variables to the Go compiler
export CGO_ENABLED:=1

#
# DOCKER FLAGS
#

export DOCKER_BUILDKIT:=1
export BUILDKIT_PROGRESS:=plain

# If set to 1, previous build caches will be ignored
DOCKER_NO_CACHE = 
DOCKER_TAG = $(VERSION)
# Alpine version used in final containers (https://hub.docker.com/_/alpine?tab=tags)
DOCKER_ALPINE_VERSION = 3.12.8
# Version of the Go compiler (https://hub.docker.com/_/golang?tab=tags)
DOCKER_GOLANG_VERSION = 1.17.1
# Alpine version used in the Golang build container (https://hub.docker.com/_/golang?tab=tags&page=1&name=alpine)
DOCKER_ALPINE_BUILD_VERSION = 3.13
DOCKERFILE_PATH = $(dir $(CURRENT_MAKEFILE))Dockerfile
DOCKER_IMAGE_BASE_PATH = docker.canopsis.net/docker/$(CANOPSIS_EDITION)

.PHONY: all docker-images deps
all: $(CMDS_BINARY_TARGETS) $(PLUGINS_BINARY_TARGETS)
docker-images: $(DOCKER_TARGETS)

#
# BINARY TARGETS
#

# "go build" does this implicitly, but doing it explicitly beforehand can
# be useful in some cases (e.g. Docker cache layer)
deps:
	@echo "==> Downloading go.mod dependencies"
	cd -- "$(ROOT_DIR)" && \
	go mod download

# Make make treat cmds as viable build targets (for single-engine build)
$(CMDS_BINARY_TARGETS):
$(PLUGINS_BINARY_TARGETS):
# Tell make where to find the binaries

# Plugin build job template
$(BUILD_DIR)/plugins/%.so: $(ROOT_DIR)/plugins/% $(ROOT_DIR)/plugins/%/*
	@echo "==> Building $@"
	cd -- "$(ROOT_DIR)" && \
	go build $(GO_BUILD_FLAGS) $(GO_BUILD_CUSTOM_FLAGS) $(LDFLAGS) \
		-buildmode=plugin \
		-o "$(abspath $@)" "$(abspath $<)"

# Build job template
$(BUILD_DIR)/cmd/%: $(ROOT_DIR)/cmd/% $(ROOT_DIR)/cmd/%/*
	@echo "==> Building $@"
	cd -- "$(ROOT_DIR)" && \
	go build $(GO_BUILD_FLAGS) $(GO_BUILD_CUSTOM_FLAGS) $(LDFLAGS) \
		-o "$(abspath $@)" "$(abspath $<)"

# To add a custom project-specific build flag
#engine-action: GO_BUILD_CUSTOM_FLAGS=--your-custom-flag

#
# DOCKER IMAGES TARGETS
#

define DOCKER_BUILD_ARGS
--platform linux/amd64 \
--build-arg "ALPINE_VERSION=${DOCKER_ALPINE_VERSION}" \
--build-arg "ALPINE_BUILD_VERSION=${DOCKER_ALPINE_BUILD_VERSION}" \
--build-arg "CANOPSIS_EDITION=${CANOPSIS_EDITION}" \
--build-arg "GOLANG_VERSION=${DOCKER_GOLANG_VERSION}"
endef

# Build a recursive .dockerignore from the monorepo (https://github.com/moby/moby/issues/20944)
$(MONOREPO_ROOT)/.dockerignore:
	cd -- "$(MONOREPO_ROOT)";\
	:> "$(MONOREPO_ROOT)/.dockerignore";\
	for file in $$(find . -name .dockerignore -type f);\
	do \
		sed -E 's;^(!?)(.*)$$;\1'$$(dirname "$$file")'/\2;' "$$file" >> .dockerignore;\
	done

.PHONY: check-docker-version docker-builder $(DOCKER_TARGETS)

# We require at least Docker 19.03 for the BuildKit features we use
check-docker-version:
	@if test "$$(docker version --format '{{.Server.Version}}' | awk -F'.' '{print $$1 $$2}')" -lt 1903 ; then\
		echo "ERROR: Docker >= 19.03 is required";\
		exit 1;\
	fi

docker-builder: $(MONOREPO_ROOT)/.dockerignore check-docker-version
	docker build \
		-t "$(DOCKER_IMAGE_BASE_PATH)/canopsis_go_bin:latest" \
		--target builder \
		$(if $(DOCKER_NO_CACHE),,--cache-from "$(DOCKER_IMAGE_BASE_PATH)/canopsis_go_bin:latest") \
		$(DOCKER_BUILD_ARGS) \
		--file "$(DOCKERFILE_PATH)" \
		"$(MONOREPO_ROOT)"

$(DOCKER_TARGETS): TARGET=final-container
$(DOCKER_TARGETS): $(MONOREPO_ROOT)/.dockerignore docker-builder
	docker build -t "$(DOCKER_IMAGE_BASE_PATH)/$(subst docker-,,$@):$(DOCKER_TAG)" \
		--build-arg CMD="$(subst docker-,,$@)" \
		--target "$(TARGET)" \
		--cache-from "$(DOCKER_IMAGE_BASE_PATH)/canopsis_go_bin:latest" \
		$(if $(DOCKER_NO_CACHE),,--cache-from "$(subst docker-,,$@):$(DOCKER_TAG)") \
		$(DOCKER_BUILD_ARGS) \
		--file "$(DOCKERFILE_PATH)" \
		"$(MONOREPO_ROOT)"

# Targets having some extra content on top of the generic final-container
docker-canopsis-api: TARGET=canopsis-api
docker-connector-junit: TARGET=connector-junit
docker-engine-webhook: TARGET=engine-webhook

#
# INSTALLATION
#

.PHONY: install

# You can modify the content of the CMDS and PLUGINS variables to change what will be built/installed
# Example : make install DESTDIR=/tmp/installroot PLUGINS=datasource/mongo.so CMDS="canopsis-api connector-junit"
# will only install the mongo datasource plugin, canopsis-api and connector-junit
install: $(dir $(CURRENT_MAKEFILE))/cmd/canopsis-reconfigure/canopsis-${CANOPSIS_EDITION}.toml $(CMDS_BINARY_TARGETS) $(PLUGINS_BINARY_TARGETS)
	# Install cmds
	install -d -- "$(DESTDIR)$(PREFIX)/opt/canopsis/bin"
	$(foreach cmd,$(CMDS),install -m 755 -- "$(BUILD_DIR)/cmd/$(cmd)" "$(DESTDIR)$(PREFIX)/opt/canopsis/bin/";)
	# Install plugins
	$(foreach plugin,$(PLUGINS),install -D -m 755 -- "$(BUILD_DIR)/plugins/$(plugin)" "$(DESTDIR)$(PREFIX)/opt/canopsis/lib/go/plugins/$(plugin)";)
	# Install configuration file
	install -d -- "$(DESTDIR)$(PREFIX)/opt/canopsis/etc"
	install -m 644 -- "$(dir $(CURRENT_MAKEFILE))/cmd/canopsis-reconfigure/canopsis-${CANOPSIS_EDITION}.toml" "$(DESTDIR)$(PREFIX)/opt/canopsis/etc/canopsis.toml"

#
# HELP
#

.PHONY: help
help:
	@echo "Available targets are:"
	@$(MAKE) -pRrq -f "$(CURRENT_MAKEFILE)" : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
	@echo
	@echo "See also: $(ROOT_DIR)/README.md"

#
# CLEANUP
#

.PHONY: clean
clean:
	rm -rf -- "$(BUILD_DIR)"
	rm -f -- "$(MONOREPO_ROOT)/.dockerignore"
