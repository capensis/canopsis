#!/usr/bin/env bash
# this script is used to deploy a fully working canopsis after the package
# has been installed manually.

set -e
set -o pipefail

case "$1" in
deploy-python|"")
	cps_engines_type=python
	;;
deploy-go)
	cps_engines_type=go
	;;
esac

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

user_home=$(su - canopsis -c 'echo -n ${HOME}')

cps_edition=core
if [ -d "${user_home}/lib/python2.7/site-packages/canopsis_cat" ]; then
	cps_edition=cat
fi

source ${user_home}/venv-ansible/bin/activate

ansible-playbook playbook/canopsis.yml \
    -e "canopsis_engines_type=$cps_engines_type" \
    -e "canopsis_edition=$cps_edition" \
    -i inventory.self
