#!/bin/bash

# usage: ./build_io.sh CONFIG_DIR

DTOOLS_DIR=$(cd $(dirname $0)/../build/depot_tools && pwd)
CONFIG_DIR=$(cd $1 && pwd)
BUILD_CONFIG_FILE=$CONFIG_DIR/CONFIG
BUILD_DIR=$(cd $(dirname $0)/../build/$(basename $CONFIG_DIR) && pwd)
GCLIENT_CONFIG_FILE=$CONFIG_DIR/GCLIENT
VERSION_CONFIG_FILE=$CONFIG_DIR/VERSION
PATCH_SCRIPT=$CONFIG_DIR/patch.sh
PATCH_DIR=$CONFIG_DIR/patch

SCRIPT_DIR=$(cd $(dirname $0) && pwd)
RTC_DIR=$BUILD_DIR/src
BUILD_AAR_CMD=$RTC_DIR/tools_webrtc/android/build_aar.py

export PATH=$DTOOLS_DIR:$PATH

source $BUILD_CONFIG_FILE
source $VERSION_CONFIG_FILE


echo "Apply patches..."

if [ -e "$PATCH_SCRIPT" ]; then
  source $PATCH_SCRIPT
fi


echo "Build Android AAR..."

BUILD_SCRIPT_OPTS="--build-dir $BUILD_DIR --build_config $CONFIG --arch $AAR_ARCH"

pushd $RTC_DIR > /dev/null

echo python $BUILD_AAR_CMD $BUILD_SCRIPT_OPTS
python $BUILD_AAR_CMD $BUILD_SCRIPT_OPTS
