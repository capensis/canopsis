#!/bin/bash
set -e
set -o pipefail

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

ansible-playbook playbook/canopsis-standalone.yml -e "canopsis_init_db=false" $@
