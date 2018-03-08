#!/bin/bash

CWD=$(dirname $(readlink -f $0))

echo "--- Cleaning old documentation"

make clean

echo "--- Generating API documentation"

(cd $CWD/developer-guide/api && ./build.sh)

echo "--- Generating HTML documentation"

make html
