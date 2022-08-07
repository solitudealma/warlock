#! /bin/bash

JS_PATH=$HOME/warlock/static/js/
JS_PATH_DIST=${JS_PATH}dist/
JS_PATH_SRC=${JS_PATH}src/

find $JS_PATH_SRC -type f -name '*.js' | sort | xargs cat > `[[ ! -d "${JS_PATH_DIST}" ]] && mkdir -p ${JS_PATH_DIST};echo ${JS_PATH_DIST}/game.js`
