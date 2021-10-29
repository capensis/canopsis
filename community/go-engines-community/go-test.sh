#!/bin/sh

[ -f Makefile ] || cd community/go-engines-community

if [ -z "$CPS_MONGO_URL" ]; then
	echo "WARNING: no CPS_MONGO_URL variable found in your environment, this will probably fail..." >&2
	sleep 5
fi

make -e test
exit $?
