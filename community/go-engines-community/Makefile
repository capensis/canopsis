ROOT_DIR := $(shell realpath --relative-to=. $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))
VPATH := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
OUTPUT_DIR=$(ROOT_DIR)/build

#List of cmds/engines names to build
CMDS := $(patsubst $(ROOT_DIR)/%,$(OUTPUT_DIR)/%,$(wildcard $(ROOT_DIR)/cmd/*))
PLUGINS := $(addsuffix .so,$(patsubst $(ROOT_DIR)/%,$(OUTPUT_DIR)/%,$(wildcard $(ROOT_DIR)/plugins/*/*)))
DOCKER_TARGETS := $(addprefix docker_,$(notdir $(wildcard $(ROOT_DIR)/cmd/*)))
INSTALL_TARGETS := $(addprefix install_cmd_,$(notdir $(wildcard $(ROOT_DIR)/cmd/*))) $(addprefix install_plugin_,$(notdir $(wildcard $(ROOT_DIR)/plugins/*/*)))


# BUILD FLAGS
BUILDDATE:=$(shell date "+%Y-%m-%dT%H:%M:%S%z")
LINKPKG="git.canopsis.net/canopsis/go-engines/lib/canopsis"
VERSION:=$(shell git for-each-ref refs/tags --sort=-authordate --format='%(refname:short)' --count=1)
GITCOMMIT:=$(shell git rev-parse HEAD)
GITBRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
#Metadata for the binaries
CUSTOM_LDFLAGS=
LDFLAGS:=-ldflags "\
	-X ${LINKPKG}.BuildGitCommit=${GITCOMMIT} \
	-X ${LINKPKG}.BuildGitBranch=${GITBRANCH} \
	-X ${LINKPKG}.BuildDate=${BUILDDATE} ${CUSTOM_LDFLAGS}"

# Use the trimpath option so that the logs do not contain the absolute paths to
# the source files.
# This makes the log lines shorter, and independent from the build machine.
GO_BUILD_FLAGS:= #-gcflags="all=-trimpath \"${GOPATH}/src\""
GO_BUILD_CUSTOM_FLAGS=

# Export configuration env variables to the go compiler
export GOOS:=linux
export GOARCH:=amd64
export CGO_ENABLED:=1
export GO111MODULE:=on

#Docker variables
DOCKER_NO_CACHE:= # If set to 1 previous build caches will be ignored
DOCKER_TAG:=$(VERSION)
DOCKER_ALPINE_VERSION:=3.13
DOCKER_GOLANG_VERSION:=1.16.4

.PHONY: build docker_images
build: $(CMDS) $(PLUGINS)
docker_images: $(DOCKER_TARGETS)

#
# BINARY TARGETS
#

# Make make treat cmds as viable build targets (for single-engine build)
$(CMDS):
$(PLUGINS):
# Tell make where to find the binaries


# Plugin build job template
$(OUTPUT_DIR)/plugins/%.so: plugins/%
	@echo "Building $@"
	@go build $(GO_BUILD_FLAGS) $(GO_BUILD_CUSTOM_FLAGS) $(LDFLAGS) \
		-buildmode=plugin \
		-o "$@" "$(abspath $^)"


# Build job template
$(OUTPUT_DIR)/cmd/%: cmd/%
	@echo "Building $@"
	@go build $(GO_BUILD_FLAGS) $(GO_BUILD_CUSTOM_FLAGS) $(LDFLAGS) \
		-o "$@" "$(abspath $^)"

# To add a custom project-specific build flag
#engine-action: GO_BUILD_CUSTOM_FLAGS=--your-custom-flag

#
# DOCKER IMAGES TARGETS
#
canopsis.toml:
	$(ROOT_DIR)/config/canopsis.toml.template core > $(ROOT_DIR)/canopsis.toml

.PHONY: docker_alpine_updated docker_builder $(DOCKER_TARGETS)
docker_alpine_updated:
	docker build \
		$(if $(DOCKER_NO_CACHE),--no-cache,--cache-from alpine_updated:latest) \
		-t alpine_updated:latest \
		--target alpine_updated \
		--build-arg "ALPINE_VERSION=${DOCKER_ALPINE_VERSION}" \
		--build-arg "GOLANG_VERSION=${DOCKER_GOLANG_VERSION}" \
		--file "$(ROOT_DIR)/Dockerfile" \

docker_builder: docker_alpine_updated
	docker build \
		-t canopsis_go_bin:latest \
		--target builder \
		--cache-from alpine_updated:latest \
		$(if $(DOCKER_NO_CACHE),,--cache-from canopsis_go_bin:latest) \
		--build-arg "ALPINE_VERSION=${DOCKER_ALPINE_VERSION}" \
		--build-arg "GOLANG_VERSION=${DOCKER_GOLANG_VERSION}" \
		--file "$(ROOT_DIR)/Dockerfile" \
		"$(ROOT_DIR)"

$(DOCKER_TARGETS): TARGET=final_container
$(DOCKER_TARGETS): docker_alpine_updated docker_builder
	docker build -t "$(subst docker_,,$@):$(DOCKER_TAG)" \
		--build-arg CMD="$(subst docker_,,$@)" \
		--target $(TARGET) \
		--cache-from alpine_updated:latest \
		--cache-from canopsis_go_bin:latest \
		$(if $(DOCKER_NO_CACHE),,--cache-from "$(subst docker_,,$@):$(DOCKER_TAG)") \
		--build-arg "ALPINE_VERSION=${DOCKER_ALPINE_VERSION}" \
		--build-arg "GOLANG_VERSION=${DOCKER_GOLANG_VERSION}" \
		--file "$(ROOT_DIR)/Dockerfile" \
		"$(ROOT_DIR)"

docker_canopsis-api: TARGET=canopsis-api

#
# Install targets
#
.PHONY: install
install: $(INSTALL_TARGETS) canopsis.toml
	install -d $(DESTDIR)$(PREFIX)/opt/canopsis/etc
	install -m 644 $(ROOT_DIR)/canopsis.toml $(DESTDIR)$(PREFIX)/opt/canopsis/etc/

# Do not make the install_* targets .PHONY as make
# does not consider implicit rule patterns for .PHONY targets
install_cmd_%: $(OUTPUT_DIR)/cmd/%
	install -d $(DESTDIR)$(PREFIX)/opt/canopsis/bin/
	install -m 755 $< $(DESTDIR)$(PREFIX)/opt/canopsis/bin/

install_plugin_%: $(OUTPUT_DIR)/plugins/%
	install -d $(DESTDIR)$(PREFIX)/opt/canopsis/lib/go/plugins/$(dir $<)
	install -m 644 $< $(DESTDIR)$(PREFIX)/opt/canopsis/lib/go/plugins/$(dir $<)/$(notdir $<)

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
	rm -f $(ROOT_DIR)/canopsis.toml
