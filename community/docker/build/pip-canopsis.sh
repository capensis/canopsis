#!/usr/bin/env bash
set -e
set -o pipefail

source ${CPS_HOME}/bin/activate

pip --no-color install --no-use-pep517 -b /tmp/pipbuild --no-deps /sources/canopsis/

rm -rf /tmp/pipbuild
rm -rf ~/.cache
