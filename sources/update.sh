#!/bin/bash

if [ `id -u` -ne 0 ]
then
	echo "You must be root..."
	exit 1
fi

SRC_PATH=`pwd`

if [ -e $SRC_PATH/canohome/lib/common.sh ]
then
	. $SRC_PATH/canohome/lib/common.sh
else
	echo "Impossible to load common.sh library"
	exit 1
fi

export_env

for pkg in $SRC_PATH/packages/*
do
	pkg=`basename $pkg`

	if [ -d $SRC_PATH/$pkg ]
	then
		echo "==> Updating package: $pkg"

		cd $SRC_PATH/$pkg

		if [ "$pkg" == "pyperfstore2" ] || [ "$pkg" == "pyperfstore3" ]
		then
			python setup.py install
		else
			echo "-- Packaging (without /etc)..."
			tar cf $SRC_PATH/$pkg.tar . --exclude=etc || exit 1

			echo "-- Extracting to $PREFIX..."
			tar xf $SRC_PATH/$pkg.tar -C $PREFIX || exit 1

			echo "-- Fix permissions..."
			for f in `tar tf $SRC_PATH/$pkg.tar`
			do
				chown $HUSER:$HGROUP $PREFIX/$f || exit 1
			done
			
			if [ -e $SRC_PATH/$pkg.tar ]
			then
				echo "-- Cleaning..."
				rm -rf $SRC_PATH/$pkg.tar
			fi
		fi
	fi
done

echo "==> Minimizing JavaScript..."
su - $HUSER -c "webcore_minimizer"		
