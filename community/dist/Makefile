VERSION = develop
BUILD_TARGET = el7 \
    el8
CANOPSIS_ARCHIVES = $(addsuffix /SOURCES/canopsis$(PACKAGE_PREFIX).tar.gz, $(addprefix build/, $(BUILD_TARGET)))
CANOPSIS_SPECS = $(addsuffix /SPECS/canopsis.spec, $(addprefix build/, $(BUILD_TARGET)))
CANOPSIS_MOCK = $(addsuffix /canopsis.cfg, $(addprefix build/, $(BUILD_TARGET)))
ARCHIVE_DIRS = community
MONOREPO_ROOT = $(abspath $(dir $(CURDIR))/../)
BUILDER_IMAGE = docker.capensis.fr/docker/images/mock-builder:v1.0.0

include $(MONOREPO_ROOT)/community/.env

TREE = build \
    $(addprefix build/, $(BUILD_TARGET)) \
    $(addsuffix /output, $(addprefix build/, $(BUILD_TARGET))) \
    $(addsuffix /SPECS, $(addprefix build/, $(BUILD_TARGET))) \
    $(addsuffix /SOURCES, $(addprefix build/, $(BUILD_TARGET)))

all: rpm

rpm: $(BUILD_TARGET)

$(BUILD_TARGET): $(CANOPSIS_ARCHIVES) $(CANOPSIS_SPECS) $(CANOPSIS_MOCK)
	docker run \
    --rm \
    -e MOCK_CONFIG=canopsis \
    -e MOCK_TARGET=SPECS \
    -v ./build/$@:/rpmbuild \
    -v ./build/$@/canopsis.cfg:/etc/mock/canopsis.cfg \
    -e DEFINE_version=${VERSION} \
    -e ENABLE_NETWORK=1 \
    --privileged=true \
    -i \
    $(BUILDER_IMAGE)

$(TREE):
	mkdir -p $@ && chmod 767 $@

$(CANOPSIS_MOCK): $(TREE)
	if [ "$(word 2,$(subst /, ,$@))" = "el7" ]; then \
	    cp centos-7-canopsis.cfg $@; \
	elif [ "$(word 2,$(subst /, ,$@))" = "el8" ]; then \
	    cp almalinux-8-canopsis.cfg $@; \
	    sed -Ei 's/nodejs-[0-9]*.[0-9]*.?[0-9]*/nodejs-$(NODE_VERSION)/g' $@; \
	fi

$(CANOPSIS_SPECS): $(TREE)
	cp canopsis$(PACKAGE_PREFIX)-$(word 2,$(subst /, ,$@)).spec $@ ; \
	sed -Ei 's/nodejs = 2:[0-9]*.[0-9]*.?[0-9]*/nodejs = 2:$(NODE_VERSION)/g' $@

$(CANOPSIS_ARCHIVES): $(TREE)
	( \
	    cd $(MONOREPO_ROOT) && \
	    git archive HEAD --prefix canopsis$(PACKAGE_PREFIX)-$(VERSION)/ $(ARCHIVE_DIRS) > \
	        $(CURDIR)/$@ \
	)

clean:
	rm -rf build

.PHONY: all rpm clean $(BUILD_TARGET)