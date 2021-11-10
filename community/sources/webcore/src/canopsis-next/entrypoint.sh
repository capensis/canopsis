#!/bin/sh
JS_STRING='window.injectedEnv = { \
  VUE_APP_API_HOST: "'"${VUE_APP_API_HOST}"'", \
};'
sed -i "s@// CONFIGURATIONS_PLACEHOLDER@${JS_STRING}@" dist/index.html
exec "$@"
