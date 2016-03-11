#!/bin/bash

PY_BIN="$PREFIX/bin/python"
INC_DIRS="/usr/include"
LOG_PATH="$BASE_PATH/log/"
INST_CONF="$SRC_PATH/build.d/"
VARLIB_PATH="$PREFIX/var/lib/pkgmgr/packages"
MESSAGES="$BASE_PATH/build_messages.txt"

CPSBUILD=1
CPSCLEAN=0
CPSNOREBUILD=0
CPSTEST=0
CPSPKG=0
CPSNODEPS=0
CPSINSTALLER=0
CPSLOGFILE=1
CPSUPDATE=0

function init_build_install() {
    export CPSBUILD
    export CPSCLEAN
    export CPSNOREBUILD
    export CPSTEST
    export CPSPKG
    export CPSNODEPS
    export CPSINSTALLER
    export CPSLOGFILE
    export CPSUPDATE

    echo "" > $MESSAGES

    mkdir -p $VARLIB_PATH >/dev/null 2>&1
    mkdir -p $LOG_PATH >/dev/null 2>&1

    rm -rf $LOG_PATH/*.log >/dev/null 2>&1
}

function check_ssl() {
    if [ -e $PREFIX/etc/ssl/ca.pem ]
    then
        openssl x509 -in $PREFIX/etc/ssl/ca.pem -noout -checkend $SSL_CHECK_SECONDS

        if [ $? -ne 0 ]
        then
            echo "-- SSL CA certificate expired (or will within $SSL_CHECK_SECONDS seconds), regenerating"
            rm $PREFIX/etc/ssl/*.pem
        fi
    fi

    if [ -e $PREFIX/etc/ssl/cert.pem ]
    then
        openssl x509 -in $PREFIX/etc/ssl/cert.pem -noout -checkend $SSL_CHECK_SECONDS

        if [ $? -ne 0 ]
        then
            echo "-- SSL certificate expired (or will within $SSL_CHECK_SECONDS seconds), regenerating"
            rm $PREFIX/etc/ssl/{cert,key,bundle}.pem
        fi
    fi
}

function purge_old_binaries() {
    if [ $CPSPKG -eq 1 ]
    then
        echo "-- Purge old binaries ..."
        rm -rf $BASE_PATH/binaries/$ARCH >/dev/null 2>&1 || true
        rm -rf $BASE_PATH/binaries/noarch >/dev/null 2>&1  || true
        rm -f $BASE_PATH/pkg.tmp >/dev/null 2>&1 || true
    fi
}

function prebuild() {
    if [ -e "$BASE_PATH/pre-build.sh" ]
    then
        . $BASE_PATH/pre-build.sh
    fi
}

function postbuild() {
    if [ -e "$BASE_PATH/post-build.sh" ]
    then
        . $BASE_PATH/post-build.sh
    fi
}

function add_message() {
    echo $@
    echo $@ >> $MESSAGES
}

function launch_log() {
    NAME=$1
    CMD=$2

    LOG=$LOG_PATH/$NAME.log

    if [ $CPSLOGFILE -eq 0 ]
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

function replace_env() {
    sed "s#@PREFIX@#$PREFIX#g" -i $PREFIX/$1
    sed "s#@HUSER@#$HUSER#g" -i $PREFIX/$1
    sed "s#@HGROUP@#$HGROUP#g" -i $PREFIX/$1
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

        replace_env etc/init.d/$1
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

        replace_env etc/$1
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

        replace_env bin/$1
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

function install_dir() {
    DIRNAME=$1
    SRCDIR=$2

    cd $SRCDIR

    echo "-- Archiving files..."
    tar cf $DIRNAME.tar $DIRNAME
    check_code $?

    echo "-- Extracting files..."
    tar mxhf $DIRNAME.tar -C $PREFIX/
    check_code $?

    echo "-- Fix permissions"
    tar tf $DIRNAME.tar | while read file
    do
        chown $HUSER:$HGROUP $PREFIX/$file
    done

    echo "-- Cleaning"
    rm $DIRNAME.tar
}

function install_pytests() {
    project=$1

    if [ -e test ]
    then
        UNITTESTDIR=$PREFIX/var/lib/canopsis/unittest/$project
        mkdir -p $UNITTESTDIR

        if [ ! -e $UNITTESTDIR/test ]
        then
            ln -s $UNITTESTDIR $UNITTESTDIR/test
        fi

        cd test
        tar -c . | tar xh -C $UNITTESTDIR
        check_code $? "Impossible to copy unittest"
        cd ..
    fi

    if [ -e features ]
    then
        FUNCTESTDIR=$PREFIX/var/lib/canopsis/functionnal-tests
        mkdir -p $FUNCTESTDIR

        cd features
        tar -c . | tar xh -C $FUNCTESTDIR
        check_code $? "Impossible to copy functionnal tests"
        cd ..
    fi
}

function extra_deps() {
    if [ $CPSNODEPS -ne 1 ]
    then
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
    fi
}

function run_clean() {
    if [ $CPSCLEAN -eq 1 ]
    then
        echo "Clean $PREFIX ..."
        echo " + kill all canopsis process ..."

        if [ -e $PREFIX/bin/hypcontrol ]
        then
            su - $HUSER -c "bash -l -c 'hypcontrol stop'"
        fi

        pkill -9 -u $HUSER
        sleep 1

        . $SRC_PATH/packages/canohome/control
        pre_remove
        post_remove
        purge
    fi
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

function build_pkg() {
    PACKAGE=$1

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

    . /$INST_CONF/$PACKAGE

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
            if [ $CPSNOREBUILD -ne 1 ]
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

        rm -f $VARLIB_PATH/$NAME.blacklist
        touch $VARLIB_PATH/$NAME.files

        cd $PREFIX

        # Find new ones
        find . -newermt "$dtstart" -type f >> $VARLIB_PATH/$NAME.files.tmp

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

        rm $VARLIB_PATH/$NAME.files.tmp

        if [ $CPSPKG -eq 1 ]
        then
            echo "-- Make package ..."

            PPATH=$BASE_PATH/binaries/$P_ARCH/$P_DIST/$P_DISTVERS

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
            echo "$NAME|$VERSION|$RELEASE|$P_DIST|$P_DISTVERS|$P_ARCH|$MD5|$REQUIRES|$DESCRIPTION" >> $BASE_PATH/pkg.tmp
        fi

        cd $SRC_PATH
    else
        echo "Impossible to build $NAME ..."
        exit 1
    fi

    add_message ""
}

function create_package_list() {
    if [ $CPSPKG -eq 1 ]
    then
        echo "-- Create packages list ..."

        cd $BASE_PATH/binaries

        echo "[" > Packages.json

        FIRST=true

        cat $BASE_PATH/pkg.tmp | while read line
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

        rm $BASE_PATH/pkg.tmp
    fi
}


function fix_permissions() {
    echo
    echo "################################"
    echo "# Fix permissions"
    echo "################################"
    echo

    mkdir -p $PREFIX/.python-eggs
    chown $HUSER:$HGROUP -R $PREFIX
    check_code $?

    echo " + Ok"
}

function launch_unittests() {
    if [ $CPSTEST -eq 1 ]
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
}


function build_installer() {
    if [ $CPSPKG -eq 1 -o $CPSINSTALLER -eq 1 ]
    then
        cd $SRC_PATH
        ./build-control/build-installer.sh
        check_code $? "Impossible to build installer"
    fi
}

function install_locales() {
    echo
    echo "################################"
    echo "# Install locales"
    echo "################################"
    echo

    cp -R $SRC_PATH/locale $PREFIX

    echo " + Ok"
}

function show_messages() {
    echo
    echo "################################"
    echo "#           MESSAGES           #"
    echo "################################"
    echo

    cat $MESSAGES
    rm $MESSAGES
}
