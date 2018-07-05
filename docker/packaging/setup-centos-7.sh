#!/bin/bash
set -e
set -o pipefail

yum makecache
yum groupinstall -y "Development tools"
yum install -y yum-utils

rm -rf /var/cache/yum
