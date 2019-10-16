#!/usr/bin/env bash
set -e
set -o pipefail

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

echo ${workdir}

src_dir="${1}"

if [ "${src_dir}" = "" ]; then
    echo "Usage: $0 <src_path>"
    echo "    <src_path> must be available inside the remote machine."
    exit 1
fi

shift

ansible-playbook -v playbook/canopsis-dev.yml --tags=cps_dev -e "canopsis_dev_src_path=${src_dir}" -e "canopsis_dev_with_deps=false" -e "canopsis_dev_without_deps=true" $@
