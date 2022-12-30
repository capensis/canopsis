#!/usr/bin/env bash
set -e
set -o pipefail

export PATH=/opt/python/3.9/bin:$PATH
source ${CPS_HOME}/bin/activate

pip3 --no-color install --no-use-pep517 --no-deps /sources/canopsis/

rm -rf ~/.cache
