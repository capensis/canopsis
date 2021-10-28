ifneq (3.81,$(firstword $(sort $(MAKE_VERSION) 3.81)))
$(error This makefile requires GNU Make >= 3.81)
endif

# Either community or pro
CANOPSIS_EDITION = community

COMMUNITY_MAKEFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
COMMUNITY_ROOT_DIR = $(realpath $(dir $(COMMUNITY_MAKEFILE)))
COMMUNITY_BUILD_OUTPUT_DIR = $(COMMUNITY_ROOT_DIR)/build

CURRENT_MAKEFILE := $(abspath $(firstword $(MAKEFILE_LIST)))
# Directory containing the main Makefile (the Makefile that was executed)
CURRENT_ROOT_DIR = $(realpath $(dir $(CURRENT_MAKEFILE)))
# Where to put the built binaries
CURRENT_BUILD_OUTPUT_DIR = $(CURRENT_ROOT_DIR)/build

# Path of the monorepo
MONOREPO_ROOT := $(realpath $(dir $(COMMUNITY_MAKEFILE))/../..)
VPATH = $(dir $(CURRENT_MAKEFILE))

# List of cmds/engines to build and install (eg: canopsis-api connector-junit...)
CMDS = $(notdir $(wildcard $(CURRENT_ROOT_DIR)/cmd/*))
# List of cmd binary paths (eg: build/cmd/canopsis-api build/cmd/connector-junit...)
CMDS_BINARY_TARGETS = $(addprefix $(CURRENT_BUILD_OUTPUT_DIR)/cmd/,$(CMDS))

# List of plugins to build and install (eg: datasource/mongo.so)
PLUGINS = $(addsuffix .so,$(patsubst $(CURRENT_ROOT_DIR)/plugins/%,%,$(wildcard $(CURRENT_ROOT_DIR)/plugins/*/*)))
PLUGINS_BINARY_TARGETS = $(addprefix $(CURRENT_BUILD_OUTPUT_DIR)/plugins/,$(PLUGINS))

# List of Docker containers targets
DOCKER_CMD_TARGETS = $(addprefix docker-cmd-,$(CMDS))

#
# BUILD FLAGS
#

SOURCE_DATE_EPOCH:=$(shell date -u +%s)
LINKPKG=git.canopsis.net/canopsis/canopsis-$(CANOPSIS_EDITION)/$(CANOPSIS_EDITION)/go-engines-$(CANOPSIS_EDITION)/lib/canopsis
VERSION:=$(shell git -C "$(MONOREPO_ROOT)" for-each-ref refs/tags --sort=-authordate --format='%(refname:short)' --count=1 2>/dev/null || echo "0.0.0")
GITCOMMIT:=$(shell git -C "$(MONOREPO_ROOT)" log -n 1 --format=oneline --abbrev-commit --pretty="%h" 2>/dev/null || echo "unknown Git commit")
# Metadata for the binaries
CUSTOM_LDFLAGS=
# TODO: simplify this after Go 1.18 with https://utcc.utoronto.ca/~cks/space/blog/programming/GoVersionOfYourSource
LDFLAGS:=-ldflags "\
	-X '$(LINKPKG).BuildGitCommit=$(GITCOMMIT)' \
	-X '$(LINKPKG).BuildDate=$(SOURCE_DATE_EPOCH)' $(CUSTOM_LDFLAGS)"

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
DOCKER_GOLANG_VERSION = 1.17.2
# Alpine version used in the Golang build container (https://hub.docker.com/_/golang?tab=tags&page=1&name=alpine)
# (note: Alpine 3.14+ is currenly on hold; see issue canopsis/canopsis-pro#3871)
DOCKER_ALPINE_BUILD_VERSION = 3.13
DOCKERFILE_PATH = $(COMMUNITY_ROOT_DIR)/Dockerfile
DOCKER_IMAGE_BASE_PATH = docker.canopsis.net/docker/$(CANOPSIS_EDITION)

.PHONY: all docker-images deps
all: $(CMDS_BINARY_TARGETS) $(PLUGINS_BINARY_TARGETS)
.DEFAULT_GOAL:=all
docker-images: $(DOCKER_CMD_TARGETS)

#
# BINARY TARGETS
#

# "go build" does this implicitly, but doing it explicitly beforehand can
# be useful in some cases (e.g. Docker cache layer)
deps:
	@echo "==> Downloading go.mod dependencies for $(CANOPSIS_EDITION)"
	cd -- "$(CURRENT_ROOT_DIR)" && \
	go mod download

# Make make treat cmds as viable build targets (for single-engine build)
$(CMDS_BINARY_TARGETS):
$(PLUGINS_BINARY_TARGETS):
# Tell make where to find the binaries

# Plugin build job template
$(CURRENT_BUILD_OUTPUT_DIR)/plugins/%.so: $(CURRENT_ROOT_DIR)/plugins/% $(CURRENT_ROOT_DIR)/plugins/%/*
	@echo "==> Building $@ for $(CANOPSIS_EDITION)"
	cd -- "$(CURRENT_ROOT_DIR)" && \
	go build $(GO_BUILD_FLAGS) $(GO_BUILD_CUSTOM_FLAGS) $(LDFLAGS) \
		-buildmode=plugin \
		-o "$(abspath $@)" "$(abspath $<)"

# Build job template
$(CURRENT_BUILD_OUTPUT_DIR)/cmd/%: $(CURRENT_ROOT_DIR)/cmd/% $(CURRENT_ROOT_DIR)/cmd/%/*
	@echo "==> Building $@ for $(CANOPSIS_EDITION)"
	cd -- "$(CURRENT_ROOT_DIR)" && \
	go build $(GO_BUILD_FLAGS) $(GO_BUILD_CUSTOM_FLAGS) $(LDFLAGS) \
		-o "$(abspath $@)" "$(abspath $<)"

# To add a custom project-specific build flag
#engine-action: GO_BUILD_CUSTOM_FLAGS=--your-custom-flag

#
# DOCKER IMAGES TARGETS
#

define DOCKER_BUILD_ARGS
--platform linux/amd64 \
--build-arg "ALPINE_VERSION=$(DOCKER_ALPINE_VERSION)" \
--build-arg "ALPINE_BUILD_VERSION=$(DOCKER_ALPINE_BUILD_VERSION)" \
--build-arg "GOLANG_VERSION=$(DOCKER_GOLANG_VERSION)" \
--build-arg "CANOPSIS_EDITION=$(CANOPSIS_EDITION)"
endef

# Build a recursive .dockerignore from the monorepo (https://github.com/moby/moby/issues/20944)
$(MONOREPO_ROOT)/.dockerignore: FORCE
	cd -- "$(MONOREPO_ROOT)" &&\
	:> "$(MONOREPO_ROOT)/.dockerignore" &&\
	for file in $$(find . -name .dockerignore -type f);\
	do \
		sed -E 's;^(!?)(.*)$$;\1'$$(dirname "$$file")'/\2;' "$$file" >> .dockerignore;\
	done
FORCE:

.PHONY: docker-check-version docker-go-builder $(DOCKER_CMD_TARGETS)

# We require at least Docker 19.03 for the BuildKit features we use
docker-check-version:
	@if test "$$(docker version --format '{{.Server.Version}}' | awk -F'.' '{print $$1 $$2}')" -lt 1903 ; then\
		echo "ERROR: Docker >= 19.03 is required";\
		exit 1;\
	fi

docker-go-builder: $(MONOREPO_ROOT)/.dockerignore docker-check-version
	docker build \
		-t "$(DOCKER_IMAGE_BASE_PATH)/canopsis-go-builder:latest" \
		--target builder \
		$(if $(DOCKER_NO_CACHE),,--cache-from "$(DOCKER_IMAGE_BASE_PATH)/canopsis-go-builder:latest") \
		$(DOCKER_BUILD_ARGS) \
		--file "$(DOCKERFILE_PATH)" \
		"$(MONOREPO_ROOT)"

$(DOCKER_CMD_TARGETS): TARGET=final-container
$(DOCKER_CMD_TARGETS): $(MONOREPO_ROOT)/.dockerignore docker-go-builder
	docker build \
		-t "$(DOCKER_IMAGE_BASE_PATH)/$(subst docker-cmd-,,$@):$(DOCKER_TAG)" \
		--target "$(TARGET)" \
		--build-arg CMD="$(subst docker-cmd-,,$@)" \
		--cache-from "$(DOCKER_IMAGE_BASE_PATH)/canopsis-go-builder:latest" \
		$(if $(DOCKER_NO_CACHE),,--cache-from "$(subst docker-cmd-,,$@):$(DOCKER_TAG)") \
		$(DOCKER_BUILD_ARGS) \
		--file "$(DOCKERFILE_PATH)" \
		"$(MONOREPO_ROOT)"

# Targets having some extra content on top of the generic final-container
docker-cmd-canopsis-api: TARGET=canopsis-api
docker-cmd-connector-junit: TARGET=connector-junit
docker-cmd-engine-webhook: TARGET=engine-webhook

#
# INSTALLATION
#

.PHONY: install

# You can modify the content of the CMDS and PLUGINS variables to change what will be built/installed
# Example : make install DESTDIR=/tmp/installroot PLUGINS=datasource/mongo.so CMDS="canopsis-api connector-junit"
# will only install the mongo datasource plugin, canopsis-api and connector-junit
install: $(COMMUNITY_ROOT_DIR)/cmd/canopsis-reconfigure/canopsis-$(CANOPSIS_EDITION).toml $(CMDS_BINARY_TARGETS) $(PLUGINS_BINARY_TARGETS)
	@echo "==> Installing files for $(CANOPSIS_EDITION)"
	# Install cmds
	install -d -m 755 -- "$(DESTDIR)$(PREFIX)/opt/canopsis/bin"
	install -c -m 755 -- $(foreach cmd,$(CMDS),"$(CURRENT_BUILD_OUTPUT_DIR)/cmd/$(cmd)") "$(DESTDIR)$(PREFIX)/opt/canopsis/bin/"
	# Install plugins
	$(foreach plugin,$(PLUGINS),install -d -m 755 -- "$(DESTDIR)$(PREFIX)/opt/canopsis/lib/go/plugins/$(dir $(plugin))" && install -c -m 644 -- "$(CURRENT_BUILD_OUTPUT_DIR)/plugins/$(plugin)" "$(DESTDIR)$(PREFIX)/opt/canopsis/lib/go/plugins/$(dir $(plugin))")
	# Install configuration files
	install -d -m 755 -- "$(DESTDIR)$(PREFIX)/opt/canopsis/etc"
	install -c -m 644 -- "$(COMMUNITY_ROOT_DIR)/cmd/canopsis-reconfigure/canopsis-$(CANOPSIS_EDITION).toml" "$(DESTDIR)$(PREFIX)/opt/canopsis/etc/canopsis.toml"

#
# TESTS
#
.PHONY: test
test: $(COMMUNITY_BUILD_OUTPUT_DIR)/cmd/ready $(COMMUNITY_BUILD_OUTPUT_DIR)/cmd/canopsis-reconfigure
	@echo "==> Running tests"
	"$(COMMUNITY_BUILD_OUTPUT_DIR)/cmd/ready" -retries 0 -timeout 2m
	"$(COMMUNITY_BUILD_OUTPUT_DIR)/cmd/canopsis-reconfigure" -conf "$(COMMUNITY_ROOT_DIR)/cmd/canopsis-reconfigure/canopsis-$(CANOPSIS_EDITION).toml" -mongoConf "$(COMMUNITY_ROOT_DIR)/config/mongo"
	# The options -p 1 -parallel 1 ensure that the tests are not run
	# concurrently, which is necessary for tests that use MongoDB, or
	# environment variables.
	cd -- "$(CURRENT_ROOT_DIR)" && \
	go test -p 1 -parallel 1 ./...

#
# HELP
#

.PHONY: help
help:
	@echo "Available targets are:"
	@$(MAKE) -pRrq -f "$(CURRENT_MAKEFILE)" : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
	@echo
	@echo "See also: $(CURRENT_ROOT_DIR)/README.md"

#
# CLEANUP
#

.PHONY: clean
clean:
	rm -rf -- "$(CURRENT_BUILD_OUTPUT_DIR)"
	rm -f -- "$(MONOREPO_ROOT)/.dockerignore"
