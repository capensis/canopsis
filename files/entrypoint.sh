#!/bin/bash

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
  su - canopsis -c "engine-launcher -e $ENGINE_MODULE -n $ENGINE_NAME -w 1 -l info"
fi
