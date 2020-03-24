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

stack_type=core
if [ -d /opt/canopsis/lib/python2.7/site-packages/canopsis_cat ]; then
	stack_type=cat
fi

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

user_home=$(su - canopsis -c 'echo -n ${HOME}')

source ${user_home}/venv-ansible/bin/activate

ansible-playbook playbook/canopsis.yml \
    -e "canopsis_engines_type=$engines_type" \
    -e "canopsis_stack_type=$stack_type" \
    -i inventory.self
