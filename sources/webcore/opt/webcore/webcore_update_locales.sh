#!/bin/bash -e

DEFAULT_DST=~/var/www/canopsis/resources/locales/

# For dev only, add DST on cli
if [ "x$1" != "x" ]; then
	if [ -e $1 ]; then
		echo "Use: $1"
		cd $1
	else
		cd $DEFAULT_DST
	fi
else
	cd $DEFAULT_DST
fi


URL="http://repo.canopsis.org/locales/"
LOCALES="fr"

for LOCALE in $LOCALES; do
	FILE="lang-$LOCALE.po"
	echo "Update '$LOCALE':"
	if [ -e $FILE ]; then
		rm -f $FILE
	fi
	wget -q $URL/$FILE -O $FILE
	echo " + $FILE: Done"
done
