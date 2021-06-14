#!/bin/bash

SRC="/usr/local/src/canopsis"

##### Uninstall
echo "-------> Uninstall"
echo "---> Clean"
pkill -u canopsis &> /dev/null || true
sleep 1
rm -Rf /opt/canopsis || true
userdel canopsis || true

##### Install from package
echo "-------> Install from packages"

echo "---> Create User"

## Requirements
useradd -m -d /opt/canopsis -s /bin/bash canopsis

## Install bootstrap
echo "---> Install bootstrap"
su - canopsis -c "mkdir -p tmp"
su - canopsis -c "rm -Rf tmp/* &> /dev/null"
su - canopsis -c "cd tmp && wget http://repo.canopsis.org/stable/canopsis_installer.tgz"
if [ $? -ne 0 ]; then exit 1; fi

su - canopsis -c "cd tmp && tar xvf canopsis_installer.tgz"
if [ $? -ne 0 ]; then exit 1; fi

su - canopsis -c "cd tmp && cd canopsis_installer && ./install.sh"
if [ $? -ne 0 ]; then exit 1; fi

echo "---> Clean bootstrap"
su - canopsis -c "rm -Rf tmp/canopsis_installer*"

## Start install
echo "---> Start install"
su - canopsis -c "ubik install --force-yes cmaster"
if [ $? -ne 0 ]; then exit 1; fi

## Check install
echo "---> Check install"
su - canopsis -c "opt/canotools/unittest.sh"
if [ $? -ne 0 ]; then exit 1; fi

echo "---> Clean"
pkill -u canopsis &> /dev/null || true
#rm -Rf /opt/canopsis
#userdel canopsis || true

## End
echo "---> Install Ok"
