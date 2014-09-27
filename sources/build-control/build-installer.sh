#!/bin/bash

### Check user
if [ `id -u` -ne 0 ]; then
	echo "You must be root ..."
	exit 1	
fi


### Configurations
SRC_PATH=`pwd`
if [ -e $SRC_PATH/canohome/lib/common.sh ]; then
	. $SRC_PATH/canohome/lib/common.sh
else
	echo "Impossible to find common's lib ..."
	exit 1
fi

echo "################################"
echo "# Make installer"
echo "################################"
cd $SRC_PATH

rm $SRC_PATH/../binaries/canopsis_installer.tgz &> /dev/null || true

echo "Copy files ..."	
cp canohome/lib/common.sh bootstrap/
check_code $?

#echo "Configure"
#BRANCH=$(git branch | grep "*" | cut -d ' ' -f2)
#if [ "$BRANCH" == "develop" ]; then
#	sed "s#stable#daily#" bootstrap/ubik.conf
#fi

echo "Create tarball installer ..."
echo "  + Make archive"
cp -R bootstrap canopsis_installer
check_code $?
tar cfz $SRC_PATH/../binaries/canopsis_installer.tgz canopsis_installer
check_code $?

echo "  + Clean"
rm -Rf canopsis_installer
check_code $?
rm bootstrap/common.sh
check_code $?

echo "  + Done"
