
PREFIX="/opt/canopsis"
HUSER="canopsis"
HGROUP="canopsis"
ARCH=`uname -m`
SUDO="sudo -E"


function check_code()
{
	if [ $1 -ne 0 ]
	then
		echo -e "Error: Code: $1"
		echo -e "Output:\n$2"
		exit $1
	fi
}

function detect_os()
{
	echo "Linux Distribution:"
	VERSION=`cat /proc/version`
	check_code $?
	DEBIAN=`echo "$VERSION" | grep -i debian | wc -l`
	UBUNTU=`echo "$VERSION" | grep -i ubuntu | wc -l`
	REDHAT=`echo "$VERSION" | grep -i redhat | wc -l`
	CENTOS=`echo "$VERSION" | grep -i centos | wc -l`
	FC=`echo "$VERSION" | grep -i "\.fc..\." | wc -l`
	ARCHL=`[ -e /etc/arch-release ] && echo 1`
	DIST_VERS=""
	
	if [ $DEBIAN -ne 0 ]
	then
		DIST="debian"
		DIST_VERS=`cat /etc/debian_version | cut -d '.' -f1`
		echo " + $DIST $DIST_VERS"
	elif [ $UBUNTU -ne 0 ]
	then
		DIST="ubuntu"
		DIST_VERS=`lsb_release -r | cut -f2`
		echo " + $DIST $DIST_VERS"
	elif [ $REDHAT -ne 0 ]
	then
		DIST="redhat"
		DIST_VERS=`lsb_release -r | cut -f2 | cut -d '.' -f1`
		echo " + $DIST $DIST_VERS"
	elif [ $CENTOS -ne 0 ]
	then
		DIST="centos"
		DIST_VERS=`lsb_release -r | cut -f2 | cut -d '.' -f1`
		echo " + $DIST $DIST_VERS"
	elif [ $FC -ne 0 ]
	then
		DIST="fedora"
		DIST_VERS=`lsb_release -r | cut -f2 | cut -d '.' -f1`
		echo " + $DIST $DIST_VERS"
	elif [ $ARCHL -ne 0 ]
	then
		DIST="archlinux"
		DIST_VERS=`pacman -Q glibc | cut -d ' ' -f2 | cut -d '-' -f1`
		echo " + $DIST $DIST_VERS"
	else
		echo " + Impossible to find distribution ..."
		exit 1
	fi
}

function launch_cmd() {
	CHECK=$1
	shift
	MYCMD=$*

	if [ "x$MYCMD" != "x" ]
	then
		if [ "x`id -un`" == "x$HUSER" ]
		then
			bash -c "$MYCMD"
			EXCODE=$?
		
			if [ $CHECK -eq 1 ]
			then
				check_code $EXCODE "Error in command '$MYCMD'..."
			else
				return $EXCODE
			fi
		elif [ `id -u` -eq 0 ]
		then
			su - $HUSER -c ". .bash_profile && $MYCMD"
			EXCODE=$?

			if [ $CHECK -eq 1 ]
			then
				check_code $EXCODE "Error in command '$MYCMD'..."
			else
				return $EXCODE
			fi
		else
			echo "Impossible to launch command with '`id -un`' ..."
			exit 1
		fi
	fi
}

function detect_numa()
{
	local CMD=`which numactl`

	if [ "x$CMD" != "x" ]
	then
		$CMD --hardware | grep 'node' | grep 'cpus' | wc -l
	else
		echo 0
	fi
}

function pkg_options()
{
	if [ $NO_ARCH == true ]
	then
		P_ARCH="noarch"
	fi
	
	if [ $NO_DISTVERS == true ]
	then
		P_DIST="nodist"
		P_DISTVERS="novers"
	elif [ $NO_DIST == true ]
	then
		P_DIST="nodist"
	fi
}

function extract_archive()
{
	if [ ! -e $1 ]
	then
		echo "Error: Impossible to find '$1' ..."
		exit 1
	fi

	echo " + Extract '$1' ..."
	tar xf $1
	check_code $? "Extract archive failure"
}

function install_init()
{
	if [ -e "$SRC_PATH/extra/init/$1.$DIST" ]
	then
		IFILE="$SRC_PATH/extra/init/$1.$DIST"
	else
		IFILE="$SRC_PATH/extra/init/$1"
	fi

	if [ -e $IFILE ]
	then
		echo " + Install init script '$1' ..."
		cp $IFILE $PREFIX/etc/init.d/$1
		check_code $? "Copy init file into init.d failure"
		sed "s#@PREFIX@#$PREFIX#g" -i $PREFIX/etc/init.d/$1
		sed "s#@HUSER@#$HUSER#g" -i $PREFIX/etc/init.d/$1
		sed "s#@HGROUP@#$HGROUP#g" -i $PREFIX/etc/init.d/$1

		check_code $? "Sed \$PREFIX,\$HUSER and \$HGROUP in init.d failure"
	else
		echo "No specific init file for $DIST"
	fi
}

function install_conf()
{
	IFILE="$SRC_PATH/extra/conf/$1"

	if [ -e $IFILE ]
	then
		echo " + Install conf file '$1' ..."
		cp $IFILE $PREFIX/etc/$1
		check_code $? "Copy conf into etc failure"
		sed "s#@PREFIX@#$PREFIX#g" -i $PREFIX/etc/$1
		sed "s#@HUSER@#$HUSER#g" -i $PREFIX/etc/$1
		sed "s#@HGROUP@#$HGROUP#g" -i $PREFIX/etc/$1
		check_code $? "Sed \$PREFIX,\$HUSER and \$HGROUP in etc failure"
	else
		echo "Error: Impossible to find '$IFILE'"
		exit 1
	fi
}

function install_bin()
{
	IFILE="$SRC_PATH/extra/bin/$1"

	if [ -e $IFILE ]
	then
		echo " + Install bin file '$1' ..."
		cp $IFILE $PREFIX/bin/$1
		check_code $? "Copy bin into bin failure"
		sed "s#@PREFIX@#$PREFIX#g" -i $PREFIX/bin/$1
		sed "s#@HUSER@#$HUSER#g" -i $PREFIX/bin/$1
		sed "s#@HGROUP@#$HGROUP#g" -i $PREFIX/bin/$1
		check_code $? "Sed \$PREFIX,\$HUSER and \$HGROUP in bin failure"
	else
		echo "Error: Impossible to find '$IFILE'"
		exit 1
	fi
}

function install_python_daemon()
{
	DPATH=$1
	DAEMON_NAME=`basename $DPATH .py`

	rm -f $PREFIX/etc/init.d/$DAEMON_NAME &>/dev/null
	check_code $? "Remove init.d script failure"
	ln -s $PREFIX/opt/canotools/daemon $PREFIX/etc/init.d/$DAEMON_NAME
	check_code $? "Symbolic link creation of daemon script failure"
}

function make_package_archive()
{
	PNAME=$1
	PPATH=$SRC_PATH/packages/$PNAME

	echo "    + Make Package archive ..."
	cd $PREFIX &> /dev/null
	tar cfj $PPATH/files.bz2 -T $PPATH/files.lst
	check_code $? "Files archive creation failure"
	cd - &> /dev/null
	
	echo "    + Check control script ..."
	touch $PPATH/control
	chmod +x $PPATH/control

	echo "    + Make final archive ..."
	cd $SRC_PATH/packages/
	tar cf $PNAME.tar $PNAME
	check_code $? "Package archive creation failure"

	echo "    + Move to binaries directory ..."
	BPATH=$SRC_PATH/../binaries/$P_ARCH/$P_DIST/$P_DISTVERS
	mkdir -p $BPATH
	check_code $? "Create Build folder failure"
	cat /proc/version > $SRC_PATH/../binaries/build.info
	mv $PNAME.tar $BPATH/
	check_code $? "Move binaries into build folder failure"

	echo "    + Clean ..."
	rm -f $PPATH/files.bz2
	check_code $? "Remove files archive failure"
}

function update_packages_list()
{
	PNAME=$1
	PPATH=$SRC_PATH/packages/$PNAME
	echo "    + Update Packages list Db ..."
	
	PKGLIST=$SRC_PATH/../binaries/Packages.list
	touch $PKGLIST

	. $PPATH/control
	check_code $? "Source control file failure"
	
	PKGMD5=$(md5sum $SRC_PATH/../binaries/$P_ARCH/$P_DIST/$P_DISTVERS/$PNAME.tar | awk '{ print $1 }')

	sed "/^$PNAME/d" -i $PKGLIST
	echo "$PNAME|$VERSION-$RELEASE||$PKGMD5|$REQUIRES|$P_ARCH|$P_DIST|$P_DISTVERS" >> $PKGLIST
}

function files_listing()
{
	local DST=$1
	if [ "x$DST" == "x" ]; then
		echo "You must specify destination ..."
		exit 1
	fi
	echo " + Files listing in $DST ..."
	mkdir -p $PREFIX
	cd $PREFIX &> /dev/null
	find ./ -type f > $DST
	find ./ -type l >> $DST
	cd - &> /dev/null|| true
	#check_code $? "List files with find failure"
}

function make_package()
{
	PNAME=$1
		
	echo " + Make package $PNAME ..."
	PPATH=$SRC_PATH/packages/$PNAME
	FLIST=$SRC_PATH/packages/files.lst
	FLIST_TMP=$SRC_PATH/packages/files.tmp
	
	mkdir -p $PPATH

	echo "    + Purge old build ..."
	rm -f $PPATH.tar &> /dev/null

	echo "    + Make files listing ..."
	files_listing "$FLIST_TMP"

	diff $FLIST $FLIST_TMP  | grep ">" | grep -v "\.pid$" | sed 's#> ##g' > $PPATH/files.lst
	check_code $?

	if [ -f $PPATH/blacklist ]; then
		echo "    + Blacklist files in listing ..."

		## blacklist files
		for line in $(cat $PPATH/blacklist); do
			cat $PPATH/files.lst | grep -v "$line" > $PPATH/files.lst.tmp
			mv $PPATH/files.lst.tmp $PPATH/files.lst
		done
	fi

	rm $FLIST_TMP
	check_code $? 'Impossible to delete tmp files listing ...'
		
	make_package_archive "$PNAME" 
}

function install_basic_source()
{
	cd $SRC_PATH
	NAME=$1
	OPTS=$2

	if [ -e "$NAME" ]
	then
		## Install file
		(cd $NAME; tar c . $OPTS | tar x -C $PREFIX)
		check_code $?
	else
		echo "Error: Impossible to find '$NAME'"
		exit 1
	fi	
}

function extra_deps()
{
	echo "Install OS dependencies for $DIST $DIST_VERS ..."
	local DEPS_FILE="$SRC_PATH/extra/dependencies/"$DIST"_"$DIST_VERS

	if [ -e $DEPS_FILE ]
	then
		bash $DEPS_FILE
	else
		echo " + Impossible to find dependencies file ($DEPS_FILE)..." 
	fi

	check_code $? "Install extra dependencies failure"
}

function run_clean()
{
	echo "Clean $PREFIX ..."
	echo " + kill all canopsis process ..."

	if [ -e $PREFIX/opt/canotools/hypcontrol ]
	then
		su - $HUSER -c ". .bash_profile; hypcontrol stop"
		check_code $? "Run hypcontrol stop failure"
	fi

	pkill -9 -u $HUSER
	sleep 1

	. $SRC_PATH/packages/canohome/control
	pre_remove
	post_remove
	purge

	rm -f $SRC_PATH/packages/files.lst &> /dev/null
}

function export_env()
{
	echo " + Fix env vars ..."
	export PATH="$PREFIX/bin:$PATH"
	export TARGET_DIR="$PREFIX/opt/rabbitmq-server"
	export SBIN_DIR="$PREFIX/bin/"
	export MAN_DIR="$PREFIX/share/man/"
	export LD_LIBRARY_PATH=$PREFIX/lib:$LD_LIBRARY_PATH
}

function pkgondemand()
{
	PNAME=$1
	echo "Make package $PNAME ..."

	if [ -e $SRC_PATH/packages/$PNAME/files.lst ]
	then
		make_package_archive "$PNAME"
	else
		echo " + Impossible to find file.lst ..."
		exit 1
	fi
	exit 0
}

function easy_install_pylib()
{
	echo "Easy install Python Library: $1 ..."
	$PREFIX/bin/easy_install -Z --prefix=$PREFIX -H None -f $SRC_PATH/externals/python-libs $1 1>> $LOG 2>> $LOG
	check_code $? "Easy install failed ..."
}