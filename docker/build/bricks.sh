#!/bin/bash
set -e
set -o pipefail

git clone https://git.canopsis.net/canopsis-ui-bricks/brick-querybuilder.git -b ${1} /sources/bricks/brick-querybuilder

rm -rf /sources/bricks/brick-querybuilder/.git

mv /sources/bricks/* ${2}/var/www/src/canopsis/

. ${2}/bin/activate

python /sources/brickmanager enable brick-querybuilder