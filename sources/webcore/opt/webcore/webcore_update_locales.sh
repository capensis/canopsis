#!/bin/bash -e

DST="~/var/www/canopsis/resources/locales/"
URL="http://repo.canopsis.org/locales/"
LOCALES="fr"

cd $DST

for LOCALE in $LOCALES; do
	FILE="lang-$LOCALE.po"
	echo "Update '$LOCALE':"
	if [ -e $FILE ]; then
		rm $FILE
	fi
	wget -q $URL/$FILE -O $FILE
	echo " + $FILE: Done"
done
