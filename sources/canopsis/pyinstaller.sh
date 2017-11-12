#!/bin/bash
workdir=$(dirname $(readlink -e $0))

cd ${workdir}

python setup.py build install || exit 1

export PYI_BIN_NAME='engine-launcher'
export PYI_DIR_NAME='engine-launcher-dir'
export PYI_SCRIPT="${workdir}/scripts/engine-launcher"
pyinstaller -y --clean -D --additional-hooks-dir ${workdir}/pyinstaller/hooks pyinstaller/canopsis.spec || exit 1

export PYI_BIN_NAME='webserver'
export PYI_DIR_NAME='webserver-dir'
export PYI_SCRIPT="${workdir}/scripts/webserver"
pyinstaller -y --clean -D --additional-hooks-dir ${workdir}/pyinstaller/hooks pyinstaller/canopsis.spec || exit 1
