#!/bin/bash

### Check user
if [ `id -u` -ne 0 ]
then
	echo "You must be root ..."
	exit 1	
fi

# Import common tools

SRC_PATH=`pwd`

if [ -e $SRC_PATH/canohome/lib/common.sh ]
then
	. $SRC_PATH/canohome/lib/common.sh
else
	echo "Impossible to find common's lib ..."
	exit 1
fi

PY_BIN=$PREFIX"/bin/python"
INC_DIRS="/usr/include"
LOG_PATH="$SRC_PATH/log/"
INST_CONF="$SRC_PATH/build.d/"

mkdir -p $LOG_PATH

function show_help()
{
	echo "Usage : ./build-install.sh <options> [packages...]"
	echo
	echo "     Install build deps, build and install Canopsis"
	echo
	echo "Options:"
	echo "    -c		->  Uninstall"
	echo "    -n		->  Don't build sources if possible"
	echo "    -u		->  Run unittest and the end"
	echo "    -p 		->  Make packages"
	echo "    -d		->  Don't check dependencies"
	echo "    -i		->  Just build installer"
	echo "    -h, help	->  Print this help"
	exit 1
}

OPT_BUILD=1
OPT_CLEAN=0
OPT_NOBUILD=0
OPT_WUT=0
OPT_MPKG=0
OPT_DCD=0
OPT_MINSTALLER=0

while getopts "cnupdhi" opt
do
	case $opt in
		c) OPT_CLEAN=1 ;;
		n) OPT_NOBUILD=1 ;;
		u) OPT_WUT=1 ;;
		p) OPT_MPKG=1 ;;
		i) OPT_MINSTALLER=1; OPT_BUILD=0;;
		d) OPT_DCD=1;;
		h) show_help ;;
		\?)
			echo "Invalid option: -$OPTARG" >&2
			show_help
		;;
	esac
done

if [ $OPT_CLEAN -eq 1 ]; then
	run_clean
	if [ "x$ARG1" == "x-c" ]; then
		exit 0
	fi
fi

export_env
detect_os

if [ $OPT_BUILD -eq 1 ]
then
	if [ $OPT_DCD -ne 1 ]
	then
		extra_deps
	fi
	
	if [ $OPT_MPKG -eq 1 ]
	then
		echo "Purge old binaries ..."
		rm -R $SRC_PATH/../binaries/$ARCH || true
		rm -R $SRC_PATH/../binaries/noarch || true
	fi

	VARLIB_PATH="$PREFIX/var/lib/pkgmgr/packages"
	mkdir -p $VARLIB_PATH
	touch $PREFIX/var/lib/pkgmgr/local_db
	
	######################################
	#  Build all packages
	######################################

	# Pre-build
	if [ -e "$SRC_PATH/pre-build.sh" ]; then
		. $SRC_PATH/pre-build.sh
	fi

	ITEMS=`ls -1 $INST_CONF | grep ".install$"`
	
	for ITEM in $ITEMS; do
		cd $SRC_PATH
	
		export MAKEFLAGS="-j$((`cat /proc/cpuinfo  | grep processor | wc -l` + 1))"
	
		NAME="x"
		VERSION="0.1"
		RELEASE="0"
		FCHECK="/tmp/notexist"
		P_ARCH=$ARCH
		P_DIST=$DIST
		P_DISTVERS=$DIST_VERS
	
		NO_ARCH=false
		NO_DIST=false
		NO_DISTVERS=false
	
		function pre_install(){	true; }
		function post_install(){ true; }
	
		. /$INST_CONF/$ITEM

		if [ "$NAME" != 'x' ]
		then
			## Check package sources
			if [ -e packages/$NAME/control ]
			then
				. packages/$NAME/control
			else
				mkdir -p packages/$NAME
				cp pkgmgr/lib/pkgmgr/control.tpl packages/$NAME/control
				sed "s#@NAME@#$NAME#g" -i packages/$NAME/control
				sed "s#@VERSION@#$VERSION#g" -i packages/$NAME/control
				sed "s#@RELEASE@#$RELEASE#g" -i packages/$NAME/control
				. packages/$NAME/control
			fi
	
			pkg_options
	
			function install(){ true; }
			function build(){ true; }
	
			. /$INST_CONF/$ITEM
	
			echo "################################"
			echo "# $NAME $VERSION"
			echo "################################"	
	
			## Build and install
			if [ ! -e $FCHECK ]
			then
				if [ $OPT_NOBUILD -ne 1 ]
				then
					echo " + Build ..."
					build
					check_code $? "Build failure"
				fi
	
				if [ $OPT_MPKG -eq 1 ]
				then
					files_listing "$SRC_PATH/packages/files.lst"
				fi
		
				echo " + Pre-install ..."	
				pre_install
	
				echo " + Install ..."
				install
				check_code $? "Install failure"
	
				echo " + Post-install ..."
				post_install
				
				if [ $OPT_MPKG -eq 1 ]
				then
					make_package $NAME
					check_code $? "Make package failure"
				fi
			else
				echo " + Allready install"
			fi
		else
			echo "Impossible to build $NAME ..."
			exit 1
		fi
	done
	
	echo "################################"
	echo "# Fix permissions"
	echo "################################"
	mkdir -p $PREFIX/.python-eggs
	chown $HUSER:$HGROUP -R $PREFIX
	check_code $?
	echo " + Ok"

	# Post-build
	if [ -e "$SRC_PATH/post-build.sh" ]
	then
		. $SRC_PATH/post-build.sh
	fi
	
	if [ $OPT_WUT -eq 1 ]
	then
		echo "################################"
		echo "# Launch Unit Tests"
		echo "################################"
		cd $SRC_PATH
		echo
		echo "Unit tests ..."
		LOG=$PREFIX/var/log/unittest.log
		launch_cmd 0 $PREFIX/opt/canotools/unittest.sh 2> $LOG 1> $LOG
		EXCODE=$?
		cp $LOG $SRC_PATH/log

		if [ $EXCODE -ne 0 ]
		then
			cat $LOG
		fi

		check_code $EXCODE "Unit tests failed ..."
		echo " + Ok"
	fi
fi

if [ $OPT_MPKG -eq 1 ] || [ $OPT_MINSTALLER -eq 1 ]
then
	cd $SRC_PATH
	./build-installer.sh
	check_code $? "Impossible to build installer"
fi
