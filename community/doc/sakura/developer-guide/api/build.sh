#!/bin/bash

CWD=$(dirname $(readlink -f $0))

echo "-- Generating RST documentation"

cat > $CWD/index.rst << "EOF"
.. _dev-api:

Canopsis API Documentation
==========================

.. toctree::
   :maxdepth: 1

EOF

cd $HOME/lib/python2.7/site-packages
for pyproj in $(ls)
do
    if [[ $pyproj == canopsis.* ]]
    then
        projectname=$(echo $pyproj | cut -d. -f2)
        projectname=$(echo $projectname | cut -d- -f1)

        echo " + $projectname..."
        echo "   $projectname/modules" >> $CWD/index.rst

        rm -rf $CWD/$projectname > /dev/null 2>&1
        sphinx-apidoc -o $CWD/$projectname $pyproj > /dev/null 2>&1
    fi
done

echo >> $CWD/index.rst
