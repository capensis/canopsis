#!/usr/bin/env bash
set -e
set -o pipefail
set -u

yum --color=never makecache
yum --color=never groupinstall -y "Development tools"
yum --color=never install -y yum-utils

rm -rf /var/cache/yum
