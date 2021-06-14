#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: before_install.sh
# Time-stamp: <2018-02-15 15:54:06>
# Copyright (C) 2018 Sergei Antipov
# Description:

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

# Latest Ansible install
pip install docker-py ansible

# Pull docker image or build it
if [ -f tests/Dockerfile.${DISTRIBUTION}_${DIST_VERSION} ]
then
    docker build --rm=true --file=tests/Dockerfile.${DISTRIBUTION}_${DIST_VERSION} --tag ${DISTRIBUTION}:${DIST_VERSION} tests
else
    docker pull ${DISTRIBUTION}:${DIST_VERSION}
fi

ln -s ${PWD} tests/greendayonfire.mongodb
