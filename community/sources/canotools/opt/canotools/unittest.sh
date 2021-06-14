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

echo -n "Wait RabbitMQ ..."
STATE=1
TRY=0
while [ $STATE -ne 0 ]; do
	if [ $TRY -eq 30 ]; then
		break
	fi
	sleep 1
	rabbitmqadmin -H 127.0.0.1 list nodes &> /dev/null
	STATE=$?
	TRY=$((TRY + 1))
	echo -n "."
done
echo
check_code $STATE "Failed to join RabbitMQ ..."
sleep 1

cd $HOME

mkdir -p var/log/unittest

UNITTESTS=`find ./ | grep Myunittest.py`
UNITTESTS="$UNITTESTS $HOME/opt/canotools/functional-test.py"

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

echo "Check celery and tasks"
runtask task_backup mongo
bfile="$HOME/var/backups/backup_mongodb.tar.gz"
if [ ! -e $bfile ]; then
	echo " + Error, backup file not found ($bfile)"
	hypcontrol stop
	exit 1
else
	echo " + Ok"
fi
hypcontrol stop
