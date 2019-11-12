#!/usr/bin/bash
set -e
set -o pipefail

source ${CPS_HOME}/bin/activate

pip install -b /tmp/pipbuild --no-deps /sources/canopsis/

rm -rf /tmp/pipbuild
rm -rf ~/.cache
