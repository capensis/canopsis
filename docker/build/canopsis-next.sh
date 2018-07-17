#!/usr/bin/env bash
set -e
set -o pipefail
set -u

cd /canopsis-next

npm install
npm run build

cp -ar dist/* /dist/
