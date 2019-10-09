#!/usr/bin/env bash
# this script is used to deploy a fully working canopsis after the package
# has been installed manually.

set -e
set -o pipefail

case "$1" in
deploy-python|"")
	engines_type=python
	;;
deploy-go)
	engines_type=go
	;;
esac

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

user_home=$(su - canopsis -c 'echo -n ${HOME}')

source ${user_home}/venv-ansible/bin/activate

ansible-playbook playbook/canopsis.yml \
    -e "canopsis_engines_type=$engines_type" \
    -i inventory.self
