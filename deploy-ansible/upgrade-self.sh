#!/usr/bin/env bash
# this script is used to deploy a fully working canopsis after the package
# has been installed manually.

set -e
set -o pipefail

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

user_home=$(su - canopsis -c 'echo -n ${HOME}')

source ${user_home}/venv-ansible/bin/activate

ansible-playbook playbook/canopsis-standalone-upgrade.yml \
    -e "canopsis_init_db=true" \
    -i inventory.self \
    --skip-tags=cps_install_package
