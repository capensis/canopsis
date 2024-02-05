VERSION = HEAD
VERSION_SAFE ?= $(shell echo "$(VERSION)" | sed -e 's/\(.*\)-\([^g][^-]*\)/\1~\2/g' -e 's/-/_/g')
BUILD_TARGET = almalinux-8-x86_64 \
	almalinux-9-x86_64
CANOPSIS_ARCHIVES = $(addsuffix /SOURCES/canopsis$(PACKAGE_PREFIX).tar.gz, $(addprefix build/, $(BUILD_TARGET)))
CANOPSIS_SPECS = $(addsuffix /SPECS/canopsis.spec, $(addprefix build/, $(BUILD_TARGET)))
CANOPSIS_MOCK = $(addsuffix /canopsis.cfg, $(addprefix build/, $(BUILD_TARGET)))
ARCHIVE_DIRS = community
MONOREPO_ROOT = $(abspath $(dir $(CURDIR))/../)
BUILDER_IMAGE = docker.capensis.fr/docker/images/mock-builder:latest

TREE = build \
    $(addprefix build/, $(BUILD_TARGET)) \
    $(addsuffix /output, $(addprefix build/, $(BUILD_TARGET))) \
    $(addsuffix /SPECS, $(addprefix build/, $(BUILD_TARGET))) \
    $(addsuffix /SOURCES, $(addprefix build/, $(BUILD_TARGET)))

all: rpm

rpm: $(BUILD_TARGET)

prep-rpm: $(CANOPSIS_ARCHIVES) $(CANOPSIS_SPECS) $(CANOPSIS_MOCK)

$(BUILD_TARGET): prep-rpm
	docker run \
	--rm \
	-e MOCK_CONFIG=canopsis \
	-e MOCK_TARGET=SPECS \
	-v ./build/$@:/rpmbuild \
	-v ./build/$@/canopsis.cfg:/etc/mock/canopsis.cfg \
  -e DEFINE_Version_safe=${VERSION_SAFE} \
	-e DEFINE_version=${VERSION} \
	-e ENABLE_NETWORK=1 \
	--privileged=true \
	-i \
	$(BUILDER_IMAGE)

$(TREE):
	mkdir -p $@

$(CANOPSIS_MOCK): $(TREE)
	if [ "$(word 2,$(subst /, ,$@))" = "almalinux-8-x86_64" ]; then \
	    cp almalinux-8-canopsis.cfg $@; \
	elif [ "$(word 2,$(subst /, ,$@))" = "almalinux-9-x86_64" ]; then \
	    cp almalinux-9-canopsis.cfg $@; \
	fi

$(CANOPSIS_SPECS): $(TREE)
	cp canopsis$(PACKAGE_PREFIX)-$(word 2,$(subst /, ,$@)).spec $@

$(CANOPSIS_ARCHIVES): $(TREE)
	( \
	    cd $(MONOREPO_ROOT) && \
	    git archive $(VERSION) --prefix canopsis$(PACKAGE_PREFIX)-$(VERSION)/ $(ARCHIVE_DIRS) > \
	        $(CURDIR)/$@ \
	)

clean:
	rm -rf build

.PHONY: all rpm clean $(BUILD_TARGET)
