#!/bin/bash

SRC_PATH=`pwd`
if [ -e ~/lib/common.sh ]; then
    . ~/lib/common.sh
else
    echo "Impossible to find common's lib ..."
    exit 1
fi

export HOME=$PREFIX
#export GUNICORN_WORKER=1

hypcontrol stop
hypcontrol start
sleep 1

    
STATE=0
TRY=0
while [ $STATE -eq 0 ]
do
    if [ $TRY -eq 30 ]
    then
        break
    fi

    sleep 1
    STATE=`launch_cmd 0 service rabbitmq-server status | grep RUNNING | wc -l`
    TRY=$((TRY + 1))
    echo -n "."
done
echo

if [ $STATE -eq 0 ]
then
    check_code 1 "Failed to test rabbit service ..."
fi
sleep 1

cd $HOME

mkdir -p var/log/unittest

UNITTESTS=`find ./ | grep Myunittest.py`
UNITTESTS="$UNITTESTS $HOME/bin/functional-test.py"

for UNITTEST in $UNITTESTS; do
    echo -n "- Proceed to $UNITTEST... " 
    python $UNITTEST > var/log/unittest/$(basename ${UNITTEST%-Myunittest.py}).log 2>&1
    EXCODE=$?

    if [ $EXCODE -eq 0 ]
    then
        echo "OK"
    else
        echo "FAIL"
    fi
done 

hypcontrol stop
