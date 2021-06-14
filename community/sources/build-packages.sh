#!/bin/bash

SRC="/usr/local/src/canopsis"
REPO_GIT="https://github.com/capensis/canopsis.git"
CMD_INSTALL="ubik install --force-yes cmaster"

BRANCH="freeze"
if [ "x$1" != "x" ]; then
	BRANCH=$1
fi

#### Git Pull
echo "-------> Clone repository"
cd /usr/local/src
if [ ! -e canopsis ]; then
	git clone $REPO_GIT
fi
echo " + Ok"

##### Build
echo "-------> Pull repository ($BRANCH)"
cd $SRC

git checkout $BRANCH
git pull origin $BRANCH

cd sources

echo "-------> Start Build"
if [ -e $SRC/builded ]; then
	./build-install.sh -cupnd
	if [ $? -ne 0 ]; then exit 1; fi
else
	./build-install.sh -cup
	if [ $? -ne 0 ]; then exit 1; fi
	touch $SRC/builded	
fi


##### Uninstall
echo "-------> Uninstall"
cd $SRC/sources
pkill -u canopsis &> /dev/null || true
./build-install.sh -c

##### Install from package
echo "-------> Install from packages"

echo "---> Install Ubik"
pip install --upgrade git+https://github.com/socketubs/ubik.git@0.1
if [ $? -ne 0 ]; then exit 1; fi

echo "---> Install Ubik tools"
pip install --upgrade git+https://github.com/capensis/ubik-toolbelt.git
if [ $? -ne 0 ]; then exit 1; fi

echo "---> Make repo"
cd $SRC/binaries
rm stable &> /dev/null
touch .repo_root
ln -s . stable
if [ $? -ne 0 ]; then exit 1; fi
ubik-repo generate stable
if [ $? -ne 0 ]; then exit 1; fi

echo "---> Start HTTP Repo"
python -m SimpleHTTPServer 8081 &
WWWCODE=$?
WWWPID=$!
if [ $WWWCODE -ne 0 ]; then exit 1; fi

echo "-----> + PID: $WWWPID"
sleep 2

echo "---> Create User"

## Requirements
useradd -m -d /opt/canopsis -s /bin/bash canopsis

## Install bootstrap
echo "---> Install bootstrap"
su - canopsis -c "mkdir -p tmp"
su - canopsis -c "rm -Rf tmp/* &> /dev/null"
su - canopsis -c "cd tmp && wget http://localhost:8081/canopsis_installer.tgz"
if [ $? -ne 0 ]; then kill -9 $WWWPID; exit 1; fi

su - canopsis -c "cd tmp && tar xvf canopsis_installer.tgz"
if [ $? -ne 0 ]; then kill -9 $WWWPID; exit 1; fi

## Configure pkgmgr
echo "---> Configure pkgmgr"
sed -i 's#repo.canopsis.org:80#localhost:8081#g' /opt/canopsis/tmp/canopsis_installer/ubik.conf

su - canopsis -c "cd tmp && cd canopsis_installer && ./install.sh"
if [ $? -ne 0 ]; then kill -9 $WWWPID; exit 1; fi

## Start install
echo "---> Start install ($CMD_INSTALL)"
su - canopsis -c "$CMD_INSTALL"
if [ $? -ne 0 ]; then kill -9 $WWWPID; exit 1; fi

## Check install
echo "---> Check install"
su - canopsis -c "opt/canotools/unittest.sh"
if [ $? -ne 0 ]; then kill -9 $WWWPID; exit 1; fi

echo "---> Clean"
kill -9 $WWWPID
pkill -u canopsis &> /dev/null || true
cd $SRC/sources
./build-install.sh -c

## End
echo "---> Package ready"
echo "+ $SRC/binaries"
