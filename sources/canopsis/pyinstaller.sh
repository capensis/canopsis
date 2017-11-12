#!/bin/bash
workdir=$(dirname $(readlink -e $0))

cd ${workdir}

if [ -z "${VIRTUAL_ENV}" ]; then
    echo "Must be in a virtualenv"
    exit 1
fi

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

python setup.py build install || exit 1

export PYI_BIN_NAME='engine-launcher'
export PYI_DIR_NAME='engine-launcher-dir'
export PYI_SCRIPT="${workdir}/scripts/engine-launcher"
pyinstaller -y --clean -D --additional-hooks-dir ${workdir}/pyinstaller/hooks pyinstaller/canopsis.spec || exit 1

export PYI_BIN_NAME='webserver'
export PYI_DIR_NAME='webserver-dir'
export PYI_SCRIPT="${workdir}/scripts/webserver"
pyinstaller -y --clean -D --additional-hooks-dir ${workdir}/pyinstaller/hooks pyinstaller/canopsis.spec || exit 1
