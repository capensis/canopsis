#!/bin/sh

[ -f Makefile ] || cd community/go-engines-community

make test
exit $?
