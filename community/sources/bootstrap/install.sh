#!/bin/bash

### Configurations
BPATH=$(pwd)
LOG="$BPATH/install.log"

if [ $(id -u) -eq 0 ]; then
	echo "Impossible to install with 'root' user ..."
	exit 1
fi

if [ -e common.sh ]; then
    . common.sh
else
	echo "Impossible to find common's lib ..."
	exit 1
fi

if [ ! -e ubik.conf ]; then
	echo "Impossible to find ubik configuration ...."
	exit 1
fi


echo
echo "#========================#"
echo "|   Canopsis Installer   |"
echo "#========================#"
echo

cd $HOME
echo > $LOG

detect_os
echo

echo " + Make directories and init environement ..."
mkdir -p etc lib var/log &>> $LOG
check_code $? "Impossible to make directories (check log: $LOG)"

cp $BPATH/common.sh lib/ &>> $LOG
check_code $? "Impossible to init environement (check log: $LOG)"

echo "   - Ok"

echo " + Configure Ubik ..."
export UBIK_CONF=$HOME/etc/ubik.conf
export PATH=$PATH:$HOME/bin
cp $BPATH/ubik.conf etc/ &>> $LOG

if [ "x$1" != "x" ]; then
	sed -i "s#stable#$1#g" $HOME/etc/ubik.conf
fi

check_code $? "Impossible to configure Ubik (check log: $LOG)"
echo "   - Ok"

echo " + Check Ubik install (package manager) ..."
ubik -v &> /dev/null
check_code $? "Impossible to find Ubik, please install Ubik with root user: 'pip install --upgrade git+https://github.com/socketubs/ubik.git@0.1'"
echo "   - Ok"

echo " + Install Canohome from canopsis package ..."
ubik install --force canohome
check_code $? "Impossible to install packages"
. .bash_profile
echo "   - Ok"
echo

echo " + Link Ubik with your Canopsis environment ..."
ubik_path=$(which ubik)
check_code $? "Impossible to find Ubik path"
if [ ! -e $ubik_path ]; then
    echo "Impossible to find Ubik path"
    exit 1
fi
ln -s $ubik_path $HOME/bin/ubik
check_code $? "Impossible to make symlink"
echo "   - Ok"

echo
echo " :: Canopsis installed"
echo

echo '   ***  /!\  Please re-login for re-load shell environement  /!\ ***'
echo
