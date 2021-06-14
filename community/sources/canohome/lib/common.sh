
PREFIX="/opt/canopsis"
HUSER="canopsis"
HGROUP="canopsis"
ARCH=`uname -m`
SUDO="sudo -E"

function check_code() {
	if [ $1 -ne 0 ]
    then
		echo -e "Error: Code: $1"
		echo -e "Output:\n$2"
		exit $1
	fi
}

function detect_os() {
	echo "Linux Distribution:"
	DIST=`python -c "import platform; print platform.dist()[0].lower().split(' ')[0]"`
	DIST_VERS=`python -c "import platform; print platform.dist()[1]"`
	echo "Dist found"
	echo $DIST
	echo $DIST_VERS
}

function launch_cmd() {
    CHECK=$1
    shift
    MYCMD=$*

    if [ "x$MYCMD" != "x" ];
    then
        if [ "x`id -un`" == "x$HUSER" ];
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

function detect_numa() {
	local CMD=`which numactl`

	if [ "x$CMD" != "x" ]
    then
		$CMD --hardware | grep 'node' | grep 'cpus' | wc -l
	else
		echo 0
	fi
}

function safe_prompt() {
    if [ "$3" == "show" ]
    then
        OPTS="-p"
    else
        OPTS="-s -p"
    fi

    empty=true

    while $empty
    do
        empty=false

        eval "read $OPTS \"$1\" $2; echo; [[ \"x\$$2\" == \"x\" ]]"

        if [ $? -eq 0 ]
        then
            empty=true

            echo "You can't let that field empty!"
        fi
    done
}

function file_append() {
    FILE=$1
    LINE=$2

    grep "$LINE" "$FILE" >/dev/null 2>/dev/null

    if [ $? -ne 0 ]
    then
        echo "$LINE" >> $FILE
    fi
}

function vercomp() {
    if [[ $1 == $2 ]]
    then
        return 0
    fi

    local IFS=.
    local i ver1=($1) ver2=($2)

    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done

    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi

        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 1
        fi
    
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    
    return 0
}
