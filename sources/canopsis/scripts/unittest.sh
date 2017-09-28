#!/bin/bash

if [ -e ~/lib/common.sh ]
then
    . ~/lib/common.sh
else
    echo "Impossible to find common's lib"
    exit 1
fi

hypcontrol start

echo -n "Wait for RabbitMQ server "
STATE=0
TRY=0
while [ $STATE -eq 0 ]
do
    if [ $TRY -eq 30 ]
    then
        break
    fi

    sleep 1
    STATE=`service rabbitmq-server status | grep -v FATAL | grep RUNNING | wc -l`
    TRY=$((TRY + 1))
    echo -n "."
done
echo

if [ $STATE -eq 0 ]
then
    check_code 1 "Failed to test RabbitMQ service ..."
fi

mkdir -p $HOME/var/log/unittest

UNITTESTDIR="$HOME/var/lib/canopsis/unittest"

UNITTESTS=`find $UNITTESTDIR -name "*.py" | grep -v "__init__.py"`
UNITTESTS="$UNITTESTS"

FAIL=0

for UNITTEST in $UNITTESTS
do
    TESTNAME=${UNITTEST#$UNITTESTDIR}

    LOGDIR=$HOME/var/log/unittest/$(dirname $TESTNAME})
    LOGFILE=$HOME/var/log/unittest/${TESTNAME%.py}.log
    mkdir -p $LOGDIR

    PYPATH=$UNITTESTDIR/$(echo $TESTNAME | cut -d/ -f2)

    echo -n "- Proceed to $TESTNAME... " 
    PYTHONPATH="$PYPATH:$PYTHONPATH" python $UNITTEST > $LOGFILE 2>&1
    EXCODE=$?

    if [ $EXCODE -eq 0 ]
    then
        echo "OK"
    else
        echo "FAIL"
        FAIL=1
    fi
done 

hypcontrol stop

if [ $FAIL -eq 1 ]
then
    exit 1
fi
