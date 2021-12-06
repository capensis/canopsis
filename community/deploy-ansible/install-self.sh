#!/usr/bin/env bash
# this script is used to deploy a fully working canopsis after the package
# has been installed manually.

set -e
set -o pipefail

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

user_home=$(su - canopsis -c 'echo -n ${HOME}')

cps_edition=community
# XXX: do a better check once Python2 is removed
if [ -d "${user_home}/lib/python2.7/site-packages/canopsis_cat" ]; then
	cps_edition=pro
fi

source ${user_home}/venv-ansible/bin/activate

ansible-playbook playbook/canopsis.yml \
    -e "canopsis_edition=$cps_edition" \
    -i inventory.self
