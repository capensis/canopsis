#!/usr/bin/env bash
set -e
set -o pipefail
set -u

cd /canopsis-next

yarn
yarn build

cp -ar dist/* /dist/
chown $FIX_OWNER -R /dist/
