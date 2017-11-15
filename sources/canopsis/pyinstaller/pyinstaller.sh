#!/bin/bash
workdir=$(dirname $(readlink -e $0))

cd ${workdir}

if [ -z "${VIRTUAL_ENV}" ]; then
    echo "Must be in a virtualenv"
    exit 1
fi

rm -rf ${workdir}/dist

export PYI_BIN_NAME='engine-launcher'
export PYI_SCRIPT="${workdir}/../scripts/engine-launcher"
pyinstaller -y --clean canopsis.spec || exit 1

export PYI_BIN_NAME='webserver'
export PYI_SCRIPT="${workdir}/../scripts/webserver"
pyinstaller -y --clean canopsis.spec || exit 1
