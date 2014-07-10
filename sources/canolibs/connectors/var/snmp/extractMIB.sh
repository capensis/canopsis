#!/bin/bash

MIB=$1

if [ "x" == "x$MIB"  ]; then
	echo "Usage: $0 <MIB>"
	exit 1
fi

MIB=`basename "$MIB"`
MIB_file=$MIB/$MIB.mib

CMIB=`echo "$MIB" | sed 's#-#_#g'`
PY_file=$MIB/$CMIB.py

echo "Extract $MIB_file in $PY_file ..."
./extractMIB.py $MIB > $PY_file
