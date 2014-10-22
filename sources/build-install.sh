#!/bin/bash

# Check user
if [ `id -u` -ne 0 ]
then
    echo "You must run this command as root ..."
    exit 1
fi


# Configurations
export SRC_PATH=`pwd`

if [ -e $SRC_PATH/canohome/lib/common.sh ]
then
    . $SRC_PATH/canohome/lib/common.sh
else
    echo "Impossible to find common's lib ..."
    exit 1
fi

if [ -e $SRC_PATH/build-control/config.sh ]
then
    . $SRC_PATH/build-control/config.sh
fi

# Check slink
if [ -e $PREFIX/.slinked ]
then
    echo "There is a slink environment installed."

    read -p "Remove slink environment ? [N/y]: " INPUT

    if [ "$(echo $INPUT | tr '[:upper:]' '[:lower:]')" == "y" ]
    then
        echo "-- Removing slink environment..."

        CPS=$PREFIX

        . $SRC_PATH/../tools/slink_utils
        remove_slinks
    else
        echo "Slink environment kept."
        exit 0
    fi
fi

PY_BIN=$PREFIX"/bin/python"
INC_DIRS="/usr/include"
LOG_PATH="$SRC_PATH/log/"
INST_CONF="$SRC_PATH/build.d/"
VARLIB_PATH="$PREFIX/var/lib/pkgmgr/packages"
FAKEROOT="$SRC_PATH/fakeroot"

MESSAGES="$SRC_PATH/build_messages.txt"
echo "" > $MESSAGES

mkdir -p $LOG_PATH
rm -rf $LOG_PATH/*.log >/dev/null 2>&1

LOGSTDOUT=0

# Functions

function add_message() {
    echo $@
    echo $@ >> $MESSAGES
}

function launch_log() {
    NAME=$1
    CMD=$2

    LOG=$LOG_PATH/$NAME.log

    if [ $LOGSTDOUT -eq 1 ]
    then
        eval $CMD
    else
        eval $CMD >>$LOG 2>&1
    fi
}

function pkg_options() {
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

function extract_archive() {
    if [ ! -e $1 ]
    then
        echo "Error: Impossible to find '$1' ..."
        exit 1
    fi

    EXTCMD=""

    if [ `echo $1 | grep 'tar.bz2$' | wc -l` -eq 1 ]
    then
        EXTCMD="tar xfj"
    elif [ `echo $1 | grep 'tar.gz$' | wc -l` -eq 1 ]
    then
        EXTCMD="tar xfz"
    elif [ `echo $1 | grep 'tgz$' | wc -l` -eq 1 ]
    then
        EXTCMD="tar xfz"
    elif [ `echo $1 | grep 'tar.xz$' | wc -l` -eq 1 ]
    then
        EXTCMD="xz -d"
    fi

    if [ "$EXTCMD" != "" ]
    then
        if [ "$EXTCMD" != "xz -d" ]
        then
            echo " + Extract '$1' ('$EXTCMD') ..."

            $EXTCMD $1
            check_code $? "Extract archive failure"
        else
            echo " + Extract '$1' ('$EXTCMD') ..."

            $EXTCMD $1
            tar xf $(echo "$1" | sed 's/.xz//')
            check_code $? "Extract archive failure"
        fi
    else
        echo "Error: Impossible to extract '$1', no command found ..."
        exit 1
    fi
}

function install_init() {
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

function install_conf() {
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

function install_bin() {
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

function update_conffiles() {
    CONFDIR=$1

    if [ -e $CONFDIR ]
    then
        echo "-- Update config files"

        find $CONFDIR -type f | while read conffile
        do
            BPATH="$SRC_PATH/build-control/etc_blacklist"
            if [ -e $BPATH ]
            then
                blacklisted=$(cat $BPATH | grep "${conffile#$CONFDIR/}" | wc -l)
            else
                blacklisted=0
            fi

            if [ $blacklisted -eq 0 ]
            then
                DEST=$PREFIX/etc/${conffile#$CONFDIR/}
                DESTDIR=$(dirname $DEST)
                copy=0

                if [ -e $DEST ]
                then
                    # Replace variables in config file for better checking
                    TMPCONF=$SRC_PATH/config.tmp
                    rm -f $TMPCONF
                    cp $conffile $TMPCONF

                    sed "s#@PREFIX@#$PREFIX#g" -i $TMPCONF
                    sed "s#@HUSER@#$HUSER#g" -i $TMPCONF
                    sed "s#@HGROUP@#$HGROUP#g" -i $TMPCONF

                    # Check the two files
                    A=$(md5sum $TMPCONF | cut -d' ' -f1)
                    B=$(md5sum $DEST | cut -d' ' -f1)

                    if [ "$A" != "$B" ]
                    then
                        if [ "$UPDATE_ETC_ACTION" == "keep" ]
                        then
                            copy=0
                        elif [ "$UPDATE_ETC_ACTION" == "replace" ]
                        then
                            copy=1
                        else
                            ask=1

                            while [ $ask -eq 1 ]
                            do
                                ask=0

                                echo "Config file: $DEST"
                                read -p "(K)eep old, (R)eplace, (V)iew diff, (E)dit (default: K): " choice <&2

                                if [ "$choice" == "R" ] || [ "$choice" == "r" ]
                                then
                                    copy=1
                                elif [ "$choice" == "V" ] || [ "$choice" == "v" ]
                                then
                                    diff $TMPCONF $DEST | more
                                    ask=1
                                elif [ "$choice" == "E" ] || [ "$choice" == "e" ]
                                then
                                    realtty="/dev/`ps | grep $$ | awk -F ' ' '{print $2}'`"

                                    if [ "x$EDITOR" == "x" ]
                                    then
                                        EDITOR="vi"
                                    fi

                                    $EDITOR $DEST < $realtty > $realtty || true
                                fi
                            done
                        fi
                    fi

                    rm $TMPCONF
                else
                    copy=1
                fi

                if [ $copy -eq 1 ]
                then
                    mkdir -p $DESTDIR
                    check_code $? "Impossible to mkdir: $DESTDIR"

                    cp -vf $conffile $DEST
                    check_code $? "Impossible to copy $conffile to $DEST"

                    chown $HUSER:$HGROUP $DEST
                    check_code $? "Impossible to chown $DEST"
                else
                    echo "Ignoring file: $conffile"
                fi
            fi
        done
    fi
}

function install_basic_source() {
    cd $SRC_PATH
    NAME=$1

    if [ -e "$NAME" ]
    then
        cd $NAME

        echo "-- Archiving files..."
        tar cf ../$NAME.tar .
        check_code $?

        echo "-- Extracting files..."
        tar mxhf ../$NAME.tar -C $PREFIX/
        check_code $?

        echo "-- Fix permissions"
        tar tf ../$NAME.tar | while read file
        do
            chown $HUSER:$HGROUP $PREFIX/$file
        done

        echo "-- Cleaning"
        rm ../$NAME.tar
    else
        echo "Error: Impossible to find '$NAME'"
        exit 1
    fi
}

function update_basic_source() {
    cd $SRC_PATH
    NAME=$1

    if [ -e "$NAME" ]
    then
        cd $NAME

        echo "-- Archiving files..."
        tar cf ../$NAME.tar . --exclude=etc
        check_code $?

        echo "-- Extracting files..."
        tar mxhf ../$NAME.tar -C $PREFIX/
        check_code $?

        echo "-- Fix permissions"
        tar tf ../$NAME.tar | while read file
        do
            chown $HUSER:$HGROUP "$PREFIX/$file"
        done

        if [ -e ./etc ]
        then
            update_conffiles ./etc
        fi

        echo "-- Cleaning"
        rm ../$NAME.tar
    else
        echo "Error: Impossible to find '$NAME'"
        exit 1
    fi
}

function extra_deps() {
    echo "Install OS dependencies for $DIST $DIST_VERS ..."

    local DEPS_FILE="$SRC_PATH/extra/dependencies/"$DIST"_"$DIST_VERS

    if [ -e $DEPS_FILE ]
    then
        bash $DEPS_FILE
    else
        VERS=$(echo $DIST_VERS | cut -d. -f1)
        DEPS_FILE="$SRC_PATH/extra/dependencies/"$DIST"_"$VERS

        if [ -e $DEPS_FILE ]
        then
            bash $DEPS_FILE
        else
            echo " + Impossible to find dependencies file ($DEPS_FILE)..."
        fi
    fi

    check_code $? "Install extra dependencies failure"
}

function run_clean() {
    echo "Clean $PREFIX ..."
    echo " + kill all canopsis process ..."

    if [ -e $PREFIX/bin/hypcontrol ]
    then
        su - $HUSER -c ". .bash_profile; hypcontrol stop"
    fi

    pkill -9 -u $HUSER
    sleep 1

    . $SRC_PATH/packages/canohome/control
    pre_remove
    post_remove
    purge
}

function export_env() {
    echo " + Fix env vars ..."

    export PATH="$PREFIX/bin:$PATH"
    export TARGET_DIR="$PREFIX/opt/rabbitmq-server"
    export SBIN_DIR="$PREFIX/bin/"
    export MAN_DIR="$PREFIX/share/man/"
    export LD_LIBRARY_PATH=$PREFIX/lib:$LD_LIBRARY_PATH
}

function easy_install_pylib() {
    echo "Easy install Python Library: $1 ..."

    $PREFIX/bin/easy_install -Z --prefix=$PREFIX -H None -f $SRC_PATH/externals/python-libs $1 1>> $LOG 2>> $LOG
    check_code $? "Easy install failed ..."
}

function show_help() {
    echo "Usage : ./build-install.sh [OPTION]"
    echo
    echo "     Build and install Canopsis dependencies and packages"
    echo
    echo "Options:"
    echo "    -g        ->  Update git submodules before build"
    echo "    -c        ->  Uninstall Canopsis"
    echo "    -n        ->  Don't build sources if possible"
    echo "    -u        ->  Run unittest at the end"
    echo "    -p        ->  Make packages"
    echo "    -d        ->  Don't check dependencies"
    echo "    -i        ->  Just build installer (for packages installation)"
    echo "    -h, help  ->  Print this help"
    exit 1
}

# Run

ARG1=$1
ARG2=$2

if [ "x$ARG1" == "xhelp" ]
then
    show_help
fi

OPT_BUILD=1
OPT_CLEAN=0
OPT_NOBUILD=0
OPT_WUT=0
OPT_MPKG=0
OPT_DCD=0
OPT_MINSTALLER=0
OPT_LOGFILE=1
OPT_GITUPDATE=0

while getopts "cnupdhilg" opt
do
    case $opt in
        c) OPT_CLEAN=1 ;;
        n) OPT_NOBUILD=1 ;;
        u) OPT_WUT=1 ;;
        p) OPT_MPKG=1 ;;
        i) OPT_MINSTALLER=1; OPT_BUILD=0;;
        l) OPT_LOGFILE=0;;
        d) OPT_DCD=1;;
        g) OPT_GITUPDATE=1;;
        h) show_help ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            show_help
        ;;
    esac
done

if [ $OPT_CLEAN -eq 1 ]
then
    run_clean

    if [ "x$ARG1" == "x-c" ]
    then
        exit 0
    fi
fi

if [ $OPT_LOGFILE -eq 0 ]
then
    LOGSTDOUT=1
fi

# Init submodules

if [ $OPT_GITUPDATE -eq 1 ]
then
    # Try to init submodule, if fail, try to use externals
    echo "-- Init submodules ..."

    cd $SRC_PATH/..

    git submodule init && git submodule update
    check_code $? "Impossible to init main submodules"

    cd $SRC_PATH/externals/
    check_code $? "Impossible to move to externals"

    git submodule init && git submodule update
    check_code $? "Impossible to init externals submodules"

    cd $SRC_PATH/..
fi

detect_os
export_env

if [ $OPT_BUILD -eq 1 ]
then
    if [ $OPT_DCD -ne 1 ]
    then
        extra_deps
    fi

    if [ $OPT_MPKG -eq 1 ]
    then
        echo "-- Purge old binaries ..."
        rm -rf $SRC_PATH/../binaries/$ARCH >/dev/null 2>&1 || true
        rm -rf $SRC_PATH/../binaries/noarch >/dev/null 2>&1  || true
        rm -f $SRC_PATH/pkg.tmp >/dev/null 2>&1 || true
    fi

    # Init package managment

    mkdir -p $VARLIB_PATH

    #  Build all packages

    # Pre-build
    if [ -e "$SRC_PATH/pre-build.sh" ]
    then
        . $SRC_PATH/pre-build.sh
    fi

    ITEMS=`ls -1 $INST_CONF | grep ".install$"`

    for ITEM in $ITEMS
    do
        dtstart=`date +"%Y%m%d %H:%M:%S.%N"`

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

        function update(){ true; }
        function install(){ true; }
        function build(){ true; }

        . /$INST_CONF/$ITEM

        if [ "$NAME" != 'x' ]
        then
            LOG=$LOG_PATH/$NAME.log
            touch $LOG

            # Check package sources
            if [ ! -e $SRC_PATH/packages/$NAME/blacklist ]
            then
                touch $SRC_PATH/packages/$NAME/blacklist
            fi

            if [ ! -e packages/$NAME/control ]
            then
                mkdir -p packages/$NAME

                function echocontrol(){
                    echo "$@" >> packages/$NAME/control
                }

                echocontrol "#!/bin/bash"
                echocontrol
                echocontrol "NAME=\"$NAME\""
                echocontrol "VERSION=\"$VERSION\""
                echocontrol "RELEASE=\"$RELEASE\""
                echocontrol "DESCRIPTION=\"$DESCRIPTION\""
                echocontrol "REQUIRES=\"$REQUIRES\""
                echocontrol
                echocontrol "function pre_install() { true; }"
                echocontrol "function post_install() { true; }"
                echocontrol "function pre_remove() { true; }"
                echocontrol "function post_remove() { true; }"
                echocontrol "function pre_update() { true; }"
                echocontrol "function post_update() { true; }"
                echocontrol "function purge() { true; }"
            fi

            chmod +x packages/$NAME/control
            . packages/$NAME/control

            pkg_options

            echo "################################################################"
            echo "#"
            add_message "# Package: $NAME v$VERSION-r$RELEASE ($DIST-$DIST_VERS $ARCH)"
            add_message "# Date: `date`"
            echo "#"
            echo "################################################################"
            add_message ""

            ## Build and install
            FORCE_UPDATE=0

            if [ $NAME == "mongodb" ]
            then
                FORCE_UPDATE=1
            fi

            ## Build and install

            #### Update package managment format
            if [ -e $VARLIB_PATH/$NAME ]
            then
                mv $VARLIB_PATH/$NAME $VARLIB_PATH/$NAME.info
            fi

            if [ $FORCE_UPDATE -eq 1 ] || [ ! -e $VARLIB_PATH/$NAME.info ]
            then
                if [ $OPT_NOBUILD -ne 1 ]
                then
                    echo "-- Build ..."
                    build
                    check_code $? "Build failure"
                fi

                echo "-- Pre-install ..."
                pre_install

                echo "-- Install ..."
                install
                check_code $? "Install failure"

                echo "-- Post-install ..."
                post_install

                echo "v${VERSION}-r${RELEASE}_${P_DIST}-${P_DISTVERS}_${P_ARCH}" > $VARLIB_PATH/$NAME.info
            else
                echo "-- Pre-update ..."
                pre_update

                echo "-- Update ..."
                update
                check_code $? "Update failure"

                echo "-- Post-update ..."
                post_update

                echo "v${VERSION}-r${RELEASE}_${P_DIST}-${P_DISTVERS}_${P_ARCH}" > $VARLIB_PATH/$NAME.info
            fi

            echo "-- Listing files ..."

            rm -f $VARLIB_PATH/$NAME.files.tmp
            rm -f $VARLIB_PATH/$NAME.blacklist
            touch $VARLIB_PATH/$NAME.files

            cd $PREFIX

            # Find new ones
            find . -newermt "$dtstart"  >> $VARLIB_PATH/$NAME.files.tmp

            # Find files matching blacklist patterns
            touch $VARLIB_PATH/$NAME.blacklist
            cat $SRC_PATH/packages/$NAME/blacklist | while read pattern
            do
                find . -wholename ".$pattern" >> $VARLIB_PATH/$NAME.blacklist
            done

            # Generate final listing
            cat $VARLIB_PATH/$NAME.files.tmp | while read file
            do
                file=${file#./}
                grep "$file" $VARLIB_PATH/$NAME.files >/dev/null 2>&1

                if [ $? -eq 1 ]
                then
                    grep "./$file" $VARLIB_PATH/$NAME.blacklist >/dev/null 2>&1

                    if [ $? -eq 1 ]
                    then
                        echo "$file" >> $VARLIB_PATH/$NAME.files
                    fi
                fi
            done

            if [ $OPT_MPKG -eq 1 ]
            then
                echo "-- Make package ..."

                PPATH=$SRC_PATH/../binaries/$P_ARCH/$P_DIST/$P_DISTVERS

                mkdir -pv $PPATH/$NAME 2> $LOG 1> $LOG

                launch_log $NAME tar cvfj $PPATH/$NAME/files.bz2 -T $VARLIB_PATH/$NAME.files
                check_code $? "files.bz2 generation failed"

                launch_log $NAME cp -v $VARLIB_PATH/$NAME.files $PPATH/$NAME/files.lst
                check_code $? "files.lst generation failed"

                launch_log $NAME cp -v $SRC_PATH/packages/$NAME/control $SRC_PATH/packages/$NAME/blacklist $PPATH/$NAME
                check_code $? "control/blacklist generation failed"

                cd $PPATH
                launch_log $NAME tar cvf $NAME.tar $NAME
                check_code $? "packaging failed"

                rm -rf $PPATH/$NAME >/dev/null 2>&1

                MD5=$(md5sum $PPATH/$NAME.tar | cut -d' ' -f1)

                echo "-- Update packages list ..."
                echo "$NAME|$VERSION|$RELEASE|$P_DIST|$P_DISTVERS|$P_ARCH|$MD5|$REQUIRES|$DESCRIPTION" >> $SRC_PATH/pkg.tmp
            fi

            cd $SRC_PATH
        else
            echo "Impossible to build $NAME ..."
            exit 1
        fi

        add_message ""
    done

    if [ $OPT_MPKG -eq 1 ]
    then
        echo "-- Create packages list ..."

        cd $SRC_PATH/../binaries

        echo "[" > Packages.json

        FIRST=true

        cat $SRC_PATH/pkg.tmp | while read line
        do
            NAME=$(echo "$line" | cut -d'|' -f1)
            VERSION=$(echo "$line" | cut -d'|' -f2)
            RELEASE=$(echo "$line" | cut -d'|' -f3)
            P_DIST=$(echo "$line" | cut -d'|' -f4)
            P_DISTVERS=$(echo "$line" | cut -d'|' -f5)
            P_ARCH=$(echo "$line" | cut -d'|' -f6)
            MD5=$(echo "$line" | cut -d'|' -f7)
            REQUIRES=$(echo "$line" | cut -d'|' -f8)
            DESCRIPTION=$(echo "$line" | cut -d'|' -f9)

            if $FIRST
            then
                FIRST=false

                echo " {" >> Packages.json
            else
                echo " },{" >> Packages.json
            fi

            echo "  \"name\": \"$NAME\"," >> Packages.json
            echo "  \"version\": \"$VERSION\"," >> Packages.json
            echo "  \"release\": \"$RELEASE\"," >> Packages.json
            echo "  \"dist\": \"$P_DIST\"," >> Packages.json
            echo "  \"vers\": \"$P_DISTVERS\"," >> Packages.json
            echo "  \"arch\": \"$P_ARCH\"," >> Packages.json
            echo "  \"md5\": \"$MD5\"," >> Packages.json

            if [ "x$REQUIRES" != "x" ]
            then
                echo "  \"requires\": [\"${REQUIRES// /\",\"}\"]," >> Packages.json
            else
                echo "  \"requires\": []," >> Packages.json
            fi

            echo "  \"description\": \"$DESCRIPTION\"" >> Packages.json
        done

        echo " }" >> Packages.json
        echo "]" >> Packages.json

        rm $SRC_PATH/pkg.tmp
    fi

    echo "################################"
    echo "# Update Python Libs"
    echo "################################"

    launch_cmd 1 update_pylibs
    echo " + Ok"

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

        launch_cmd 0 $PREFIX/bin/unittest.sh 2> $LOG 1> $LOG
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

if [ $OPT_MPKG -eq 1 -o $OPT_MINSTALLER -eq 1 ]
then
    cd $SRC_PATH
    ./build-control/build-installer.sh
    check_code $? "Impossible to build installer"
fi

echo
echo "################################"
echo "# Installing locales"
echo "################################"
echo

cp -R $SRC_PATH/locale $PREFIX

echo " + Ok"


echo
echo "################################"
echo "#           MESSAGES           #"
echo "################################"
echo

cat $MESSAGES
rm $MESSAGES

echo " -- You can now run Canopsis: hypcontrol start"
