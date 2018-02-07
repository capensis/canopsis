#!/bin/bash

# Replacment of old doc build system. Might be broken.

HUSER=canopsis
HGROUP=canopsis
PREFIX="/opt/canopsis"
SRC_PATH="/vagrant/sources"

function check_code() {
    if [ $1 -ne 0 ]
    then
        echo -e "Error: Code: $1"
        echo -e "Output:\n$2"
        exit $1
    fi
}

function pip_install() {
    echo "pip install $*"

    pip install --no-index --find-links=file://${SRC_PATH}/externals/python-libs $*
    check_code $? "Pip install failed ..."
}

pip_install "MarkupSafe==0.23"
pip_install "docutils==0.12"
pip_install "Pygments==1.6"
pip_install "jinja2==2.7.3"
pip_install "Sphinx==1.2.3"


rm -rf $PREFIX/var/lib/canopsis/doc
rm -rf $PREFIX/var/www/src/doc

mkdir -p $PREFIX/var/lib/canopsis/doc
mkdir -p $PREFIX/var/www/src/doc

echo "--- Installing documentation"

cp -r * $PREFIX/var/lib/canopsis/doc

chown -R $HUSER:$HGROUP $PREFIX/var/lib/canopsis/doc
chown -R $HUSER:$HGROUP $PREFIX/var/www/src/doc

echo "--- Building documentation"

./build.sh

#cp -r $PREFIX/var/lib/canopsis/doc/_build/html/* $PREFIX/var/www/src/doc
cp -r _build/html/* $PREFIX/var/www/src/doc
