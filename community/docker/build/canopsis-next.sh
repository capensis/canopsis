#!/usr/bin/env bash
set -e
set -o pipefail
set -u

cd /canopsis-next

export FORCE_COLOR=0 NPM_CONFIG_COLOR=false NPM_CONFIG_PROGRESS=false NPM_CONFIG_SPIN=false
yarn
NODE_ENV=production yarn build --mode production

cp -ar dist/* /dist/
chown $FIX_OWNER -R /dist/
