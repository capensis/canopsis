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

function export_env(){
	echo " + Fix env vars ..."
	export PATH="$PREFIX/bin:$PATH"
	export TARGET_DIR="$PREFIX/opt/rabbitmq-server"
	export SBIN_DIR="$PREFIX/bin/"
	export MAN_DIR="$PREFIX/share/man/"
	export LD_LIBRARY_PATH=$PREFIX/lib:$LD_LIBRARY_PATH
}

export_env

for pkg in $SRC_PATH/packages/*
do
	pkg=`basename $pkg`

	if [ -d $SRC_PATH/$pkg ]
	then
		echo "==> Updating package: $pkg"

		cd $SRC_PATH/$pkg

		if [ "$pkg" == "webcore" ]
		then
			export CPS_PREFIX="$PREFIX/etc/"
			python setup.py install --no-conf
		elif [ "$pkg" == "python" ]
		then
			PROJECTS[0]='common'
			PROJECTS[1]='configuration'
			PROJECTS[2]='timeserie'
			PROJECTS[3]='storage'
			PROJECTS[4]='context'
			PROJECTS[5]='perfdata'
			PROJECTS[6]='mongo'
			PROJECTS[7]='old'
			PROJECTS[8]='engines'
			PROJECTS[9]='connectors'
			PROJECTS[10]='tools'
			PROJECTS[11]='cli'
			PROJECTS[12]='topology'
			PROJECTS[13]='organisation'
			PROJECTS[14]='auth'

			for project in "${PROJECTS[@]}"
			do
				echo "-- Install project: $project"
				cd $SRC_PATH/$pkg/$project
				export CPS_PREFIX="$PREFIX/etc/"
				python setup.py install --no-conf || exit 1
			done

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
#su - $HUSER -c "webcore_minimizer"
