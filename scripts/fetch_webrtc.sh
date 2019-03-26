#!/bin/bash

# usage: ./fetch_webrtc.sh CONFIG_DIR

URL=https://chromium.googlesource.com/external/webrtc
DTOOLS=$(cd $(dirname $0)/../build/depot_tools && pwd)

CONFIG_DIR=$(cd $1 && pwd)
VERSION_CONFIG=$CONFIG_DIR/VERSION
GCLIENT_CONFIG=$CONFIG_DIR/GCLIENT

BUILD_DIR=$(dirname $0)/../build/$(basename $CONFIG_DIR)
mkdir -p $BUILD_DIR
BUILD_DIR=$(cd $BUILD_DIR && pwd)
RTC_DIR=$BUILD_DIR/src


export PATH=$DTOOLS:$PATH

source $VERSION_CONFIG
GCLIENT_CONFIG_SPEC=`tr -d "\n" < ${GCLIENT_CONFIG}`

if [ $BRANCH -gt 72 ]; then
  BRANCH_HEADS=m$BRANCH
else
  BRANCH_HEADS=$BRANCH
fi

echo "Checkout the code with release branch M$BRANCH ($REVISION)..."

pushd $BUILD_DIR > /dev/null

echo "Initialize gclient..."
gclient config --spec "$GCLIENT_CONFIG_SPEC"

echo "Sync..."
gclient sync --nohooks --with_branch_heads -v -R

pushd $RTC_DIR > /dev/null
git submodule foreach "'git config -f \$toplevel/.git/config submodule.\$name.ignore all'"
git config diff.ignoreSubmodules all

git fetch origin
git checkout -B "M$BRANCH" remotes/branch-heads/$BRANCH_HEADS
git checkout $REVISION

# Android だと途中でライセンスの同意が必要になるので yes を置く
yes | gclient sync --with_branch_heads -v -R

gclient runhooks -v

