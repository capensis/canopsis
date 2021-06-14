#!/bin/bash

if [ "$1" != "" ]; then
    MONGO_URL=$1
else
    echo "Please provide mongo connection url"
    exit 1
fi

BASEDIR=$(dirname "$0")

for entry in "$BASEDIR"/*.js
do
  echo "Start migration: $entry"
  mongo $MONGO_URL $entry
  echo "Finish migration: $entry"
done
