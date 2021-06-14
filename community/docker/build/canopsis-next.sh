#!/usr/bin/env bash
set -e
set -o pipefail
set -u

cd /canopsis-next

yarn
NODE_ENV=production yarn build --mode production

cp -ar dist/* /dist/
chown $FIX_OWNER -R /dist/
