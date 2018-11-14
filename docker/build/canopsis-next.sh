#!/usr/bin/env bash
set -e
set -o pipefail
set -u

cd /canopsis-next

npm install
NODE_ENV=production npm run build

cp -ar dist/* /dist/
chown $FIX_OWNER -R /dist/
