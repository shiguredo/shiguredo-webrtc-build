#!/bin/sh

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
BUILD_IOS_CMD=$RTC_DIR/tools_webrtc/ios/build_ios_libs.sh
BUILD_LIB_PATH=$BUILD_DIR/WebRTC.framework
BUILD_INFO_FILE=$BUILD_LIB_PATH/build_info.json

export PATH=$DTOOLS_DIR:$PATH

source $BUILD_CONFIG_FILE
source $VERSION_CONFIG_FILE


echo "Apply patches..."

if [ -e "$PATCH_SCRIPT" ]; then
  source $PATCH_SCRIPT
fi


echo "Build iOS framework..."

BUILD_SCRIPT_OPTS="-o $BUILD_DIR --build_config $CONFIG --arch $IOS_ARCH"
if [ "$IOS_BITCODE" = "true" ]; then
  BUILD_SCRIPT_OPTS="$BUILD_SCRIPT_OPTS --bitcode"
fi
if [ "$VP9" = "true" ]; then
  BUILD_SCRIPT_OPTS="$BUILD_SCRIPT_OPTS --vp9"
fi

echo $BUILD_IOS_CMD $BUILD_SCRIPT_OPTS
$BUILD_IOS_CMD $BUILD_SCRIPT_OPTS


echo "Generate build_info.json..."
cat <<EOF > $BUILD_INFO_FILE
{
    "webrtc_version" = "$BRANCH",
    "webrtc_revision" = "$REVISION"
}
EOF

zip -rq $BUILD_LIB_PATH.zip $BUILD_LIB_PATH 
