#!/usr/bin/env bash
set -e
set -o pipefail
set -u

workdir=$(dirname $(readlink -e $0))
source ${workdir}/../../build-env.sh

rm -rf docker/bricks && mkdir -p docker/bricks

git clone https://git.canopsis.net/canopsis-ui-bricks/brick-querybuilder.git -b ${CANOPSIS_UIV2_BRICKS_TAG} docker/bricks/brick-querybuilder
git clone https://git.canopsis.net/canopsis-ui-bricks/brick-listalarm.git -b ${CANOPSIS_UIV2_BRICKS_TAG} docker/bricks/brick-listalarm
git clone https://git.canopsis.net/canopsis-ui-bricks/brick-timeline.git -b ${CANOPSIS_UIV2_BRICKS_TAG} docker/bricks/brick-timeline

rm -rf /sources/bricks/brick-querybuilder/.git
rm -rf /sources/bricks/brick-listalarm/.git
rm -rf /sources/bricks/brick-timeline/.git
