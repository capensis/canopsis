#!/bin/bash

# Check user

if [ `id -u` -ne 0 ]
then
    echo "You must run this command as root ..."
    exit 1
fi

# Configurations

export BASE_PATH=$(dirname $(readlink -f $0))
export SRC_PATH="$BASE_PATH/sources"

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

if [ -e $SRC_PATH/build-control/functions.sh ]
then
    . $SRC_PATH/build-control/functions.sh
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

        . $BASE_PATH/tools/slink_utils
        remove_slinks
    else
        echo "Slink environment kept."
        exit 0
    fi
fi

# Functions

function init_submodules() {
    if [ $CPSUPDATE -eq 1 ]
    then
        # Try to init submodule, if fail, try to use externals
        echo "-- Init submodules ..."

        cd $BASE_PATH

        git submodule init && git submodule update
        check_code $? "Impossible to init main submodules"

        cd $SRC_PATH/externals/
        check_code $? "Impossible to move to externals"

        git submodule init && git submodule update
        check_code $? "Impossible to init externals submodules"

        cd $BASE_PATH
    fi
}

function show_help() {
    echo "Usage : ./build-install.sh [options]"
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
    echo "    -o pkg    ->  Build single package (can be used many time)"
    echo "    -h, help  ->  Print this help"
}

# Run

CPSBUILDONE=0
PACKAGES=""

if [ "x$1" == "xhelp" ]
then
    show_help
fi

while getopts "cnupdhilgo:" opt
do
    case $opt in
        c)
            CPSCLEAN=1
        ;;

        n)
            CPSNOREBUILD=1
        ;;

        u)
            CPSTEST=1
        ;;
        
        p)
            CPSPKG=1
        ;;
        
        i)
            CPSINSTALLER=1
            CPSBUILD=0
        ;;
        
        l)
            CPSLOGFILE=0
        ;;
        
        d)
            CPSNODEPS=1
        ;;
        
        g)
            CPSUPDATE=1
        ;;

        o)
            CPSBUILDONE=1
            PACKAGES="$PACKAGES $OPTARG"
        ;;
        
        h)
            show_help
            exit 0
        ;;
        
        \?)
            echo "Invalid option: -$OPTARG" >&2
            show_help
            exit 1
        ;;
    esac
done

run_clean

if [ $CPSCLEAN -eq 1 ] && [ "$1" == "-c" ]
then
    exit 0
fi

init_build_install
init_submodules
detect_os

if [ $CPSBUILD -eq 1 ]
then
    extra_deps
    export_env
    check_ssl

    if [ $CPSBUILDONE -eq 0 ]
    then
        purge_old_binaries
        prebuild

        ITEMS=$(ls -1 $INST_CONF | grep ".install$")

        for ITEM in $ITEMS
        do
            build_pkg $ITEM || break
        done

        postbuild
        create_package_list
    else
        for PACKAGE in $PACKAGES
        do
            build_pkg $(find $INST_CONF -name "*$PACKAGE.install" -printf "%f")
        done
    fi

    fix_permissions
    launch_unittests
else
    export_env
fi

build_installer
install_locales
show_messages

echo " -- You can now run Canopsis: hypcontrol start"
