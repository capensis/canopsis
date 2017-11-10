#!/bin/sh
find canopsis -type f -name "*.py"|grep -v "__init__"|grep -v "\.cli\."|grep "\.py$"|sed -r -e 's/\.\//canopsis./g' -e 's/\.py$//g' -e 's/\//./g' |sort|uniq
