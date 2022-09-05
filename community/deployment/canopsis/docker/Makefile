ifneq (3.81,$(firstword $(sort $(MAKE_VERSION) 3.81)))
$(error This makefile requires GNU Make >= 3.81)
endif

CURRENT_MAKEFILE := $(abspath $(firstword $(MAKEFILE_LIST)))
# Directory containing the main Makefile (the Makefile that was executed)
CURRENT_ROOT_DIR = $(realpath $(dir $(CURRENT_MAKEFILE)))
# Where to put the built binaries
CURRENT_BUILD_OUTPUT_DIR = $(CURRENT_ROOT_DIR)/build

# Either community or pro
CANOPSIS_EDITION = community

VERSION:=$(shell git -C "$(MONOREPO_ROOT)" for-each-ref refs/tags --sort=-authordate --format='%(refname:short)' --count=1 2>/dev/null || echo "0.0.0")

TARBALL_FILE = $(CURRENT_BUILD_OUTPUT_DIR)/canopsis-$(CANOPSIS_EDITION)-docker-compose-$(VERSION).tar

.PHONY: all
all: tarball
.DEFAULT_GOAL:=all

.PHONY: tarball
tarball: $(CURRENT_BUILD_OUTPUT_DIR) $(CURRENT_ROOT_DIR)
	echo $(CANOPSIS_EDITION)
	tar -C $(dir $(CURRENT_ROOT_DIR)) \
		--create \
		--verbose \
		--file $(TARBALL_FILE)\
		--dereference \
		--exclude-vcs --exclude-vcs-ignores \
		--exclude=Makefile --exclude=build \
		--transform='s;^$(notdir $(CURRENT_ROOT_DIR))/;;' \
		$(notdir $(CURRENT_ROOT_DIR))

$(CURRENT_BUILD_OUTPUT_DIR):
	mkdir -p $@


.PHONY: clean

clean:
	rm -rf $(CURRENT_BUILD_OUTPUT_DIR)

.PHONY: help
help:
	@echo "This makefile builds portable tarballs of the Canopsis docker-compose environment"
	@echo "Available targets are:"
	@$(MAKE) -pRrq -f "$(CURRENT_MAKEFILE)" : 2>/dev/null | awk -F':' '/^[a-zA-Z0-9][^$$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)print A[i]}' | sort -u
	@echo
	@echo "See also: $(CURRENT_ROOT_DIR)/README.md"

