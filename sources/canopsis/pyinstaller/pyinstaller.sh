#!/bin/bash
workdir=$(dirname $(readlink -e $0))

cd ${workdir}

if [ -z "${VIRTUAL_ENV}" ]; then
    echo "Must be in a virtualenv"
    exit 1
fi

function remove_hooks() {
    find ${VIRTUAL_ENV} -name "hook-_tkinter.pyc" -delete
    find ${VIRTUAL_ENV} -name "pyi_rth__tkinter.pyc" -delete
    f=$(find ${VIRTUAL_ENV} -name "hook-_tkinter.py")
    if [ ! "${f}" = "" ]; then
        echo -n "" > "${f}"
    fi
    f=$(find ${VIRTUAL_ENV} -name "pyi_rth__tkinter.py")
    if [ ! "${f}" = "" ]; then
        echo -n "" > "${f}"
    fi
}

remove_hooks

rm -rf ${workdir}/dist

export PYI_BIN_NAME='engine-launcher'
export PYI_SCRIPT="${workdir}/../scripts/engine-launcher"
pyinstaller -y --clean canopsis.spec || exit 1

export PYI_BIN_NAME='webserver'
export PYI_SCRIPT="${workdir}/../scripts/webserver"
pyinstaller -y --clean canopsis.spec || exit 1
