#!/bin/sh

# usage: ./fetch_depot_tools.sh DIR

URL=https://chromium.googlesource.com/chromium/tools/depot_tools.git
DIR=$1
DTOOLS_DIR=$DIR/depot_tools

if [ ! -e $DTOOLS_DIR ]; then
  echo "Get depot_tools..."
  git clone $URL $DTOOLS_DIR
else
	echo "Update depot_tools..."
  git -C $DTOOLS_DIR pull
fi
