#!/bin/bash

## Check deps
BIN_AVCONV=`which avconv`
BIN_PHANTOMJS=`which phantomjs`
BIN_CASPERJS=`which casperjs`

if [ "x$BIN_PHANTOMJS" == "x" ]; then
	echo "Impossible to find 'phantomjs'"
	exit 1
fi

if [ "x$BIN_CASPERJS" == "x" ]; then
	echo "Impossible to find 'casperjs'"
	exit 1
fi

## Make captures dir
mkdir -p captures
rm -f captures/* || true

## Build global test file
TFILE="tests.js"
cat main.js > $TFILE

echo "Build tests file '$TFILE'"

if [ "x$1" != "x" ]; then
	files="test.d/000_auth.js $1 test.d/999_logout.js"
else
	files=test.d/*
fi

for file in $files; do
	if [[ -f $file ]]; then
		echo "  - $file"
		echo >> $TFILE
		echo "casper.then(function() { casper.echo('\n###########################\n# $file\n###########################', 'COMMENT'); });" >> $TFILE
		cat "$file" >> $TFILE
	fi
done
echo " + Done"

## Run
$BIN_CASPERJS test $TFILE
CODE=$?

if [ "x$BIN_AVCONV" != "x" ]; then
	## Conv PNG to video
	echo
	echo "Build video ..."
	LOG="captures/avconv.log"
	#avconv -y -threads 8 -r 6 -i captures/step-%d.png -vcodec mpeg4 tests.mp4
	$BIN_AVCONV -y -r 6 -i captures/step-%d.png tests.webm 1> $LOG 2> $LOG
	echo " + Done"
fi

## Quit
echo "Exit with code: $CODE"
exit $CODE

