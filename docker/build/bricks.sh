#!/usr/bin/bash
set -e
set -o pipefail

rm -rf docker/bricks && mkdir -p docker/bricks

git clone https://git.canopsis.net/canopsis-ui-bricks/brick-querybuilder.git -b ${1} docker/bricks/brick-querybuilder

rm -rf /sources/bricks/brick-querybuilder/.git