#!/bin/sh

# usage: ./patch.sh CONFIG_DIR

DTOOLS_DIR=$(cd $(dirname $0)/../build/depot_tools && pwd)
CONFIG_DIR=$(cd $1 && pwd)
BUILD_CONFIG_FILE=$CONFIG_DIR/CONFIG
BUILD_DIR=$(dirname $0)/../build/$(basename $CONFIG_DIR)
GCLIENT_CONFIG_FILE=$CONFIG_DIR/GCLIENT
VERSION_CONFIG_FILE=$CONFIG_DIR/VERSION
SCRIPT_DIR=$(dirname $0)
RTC_DIR=$BUILD_DIR/src
BUILD_IOS_CMD=$RTC_DIR/tools_webrtc/ios/build_ios_libs.sh
BUILD_LIB_PATH=$BUILD_DIR/WebRTC.framework
BUILD_INFO_FILE=$BUILD_LIB_PATH/build_info.json

source $SCRIPT
