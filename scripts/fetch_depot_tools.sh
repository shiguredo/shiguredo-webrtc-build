#!/bin/sh

# usage: ./fetch_depot_tools.sh DIR

URL=https://chromium.googlesource.com/chromium/tools/depot_tools.git
DIR=$1/depot_tools

if [ ! -e $DIR ]; then
  echo "Get depot_tools..."
  git clone $URL $DIR
else
	echo "Update depot_tools..."
  git -C $DIR pull
fi
