CURRENT_MAKEFILE := $(abspath $(firstword $(MAKEFILE_LIST)))

# Path of the monorepo
MONOREPO_ROOT := $(realpath $(dir $(CURRENT_MAKEFILE))/../../../../)
include $(MONOREPO_ROOT)/community/.env

ifneq ($(REQUIRED_MAKE_VERSION),$(firstword $(sort $(MAKE_VERSION) $(REQUIRED_MAKE_VERSION))))
$(error This makefile requires GNU Make >= $(REQUIRED_MAKE_VERSION))
endif

# Directory containing the main Makefile (the Makefile that was executed)
CURRENT_ROOT_DIR = $(realpath $(dir $(CURRENT_MAKEFILE)))
# Where to put the built binaries
CURRENT_BUILD_OUTPUT_DIR = $(CURRENT_ROOT_DIR)/build

# Either community or pro
CANOPSIS_EDITION = community

VERSION:=$(shell (git -C "$(MONOREPO_ROOT)" tag --list --sort=-taggerdate --merged 2>/dev/null || echo "0.0.0") | head -1)
VERSION_TIME:=$(shell git for-each-ref --format="%(taggerdate:unix)" refs/tags/$(VERSION) 2> /dev/null || stat -c %Y .env)
# If returned tag match ci/\#5402 replace by develop
ifeq ($(VERSION),ci/\#5402)
    VERSION := develop
endif

TARBALL_NAME = canopsis-$(CANOPSIS_EDITION)-docker-compose-$(VERSION)
TARBALL_FILE = $(CURRENT_BUILD_OUTPUT_DIR)/$(TARBALL_NAME).tar.gz

.PHONY: all
all: tarball
.DEFAULT_GOAL:=all

.PHONY: tarball
tarball: $(CURRENT_BUILD_OUTPUT_DIR) $(CURRENT_ROOT_DIR)
	echo $(CANOPSIS_EDITION)
	sed 's/^CPS_EDITION=.*//g' "$(CURRENT_ROOT_DIR)/.env" > env
	sed -i '1iCPS_EDITION=$(CANOPSIS_EDITION)' "$(CURRENT_ROOT_DIR)/env"
	tar -z -C $(dir $(CURRENT_ROOT_DIR)) \
		--create \
		--verbose \
		--file $(TARBALL_FILE) \
		--dereference \
		--mtime @$(VERSION_TIME) \
		--sort=name \
		--owner=0 --group=0 --numeric-owner \
		--exclude-vcs \
		--exclude=Makefile --exclude=build --exclude=.env \
		--transform='s;^$(notdir $(CURRENT_ROOT_DIR))\(/\|$$\);$(TARBALL_NAME)\1;' \
		--transform='s;/env$$;/.env;' \
		$(notdir $(CURRENT_ROOT_DIR))

$(CURRENT_BUILD_OUTPUT_DIR):
	mkdir -p $@


.PHONY: clean

clean:
	rm -rf $(CURRENT_BUILD_OUTPUT_DIR)
	rm -f env

.PHONY: help
help:
	@echo "This makefile builds portable tarballs of the Canopsis docker-compose environment"
	@echo "Available targets are:"
	@$(MAKE) -pRrq -f "$(CURRENT_MAKEFILE)" : 2>/dev/null | awk -F':' '/^[a-zA-Z0-9][^$$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)print A[i]}' | sort -u
	@echo
	@echo "See also: $(CURRENT_ROOT_DIR)/README.md"


