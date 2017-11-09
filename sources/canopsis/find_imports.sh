#!/bin/sh
tree -nif --prune canopsis|grep -v "__init__"|grep "\.py$"|sed -r -e 's/\.\//canopsis./g' -e 's/\.py$//g' -e 's/\//./g' |sort|uniq
