#!/bin/sh

patch -buN $RTC_DIR/tools_webrtc/ios/build_ios_libs.py $PATCH_DIR/build_ios_libs.py.diff
#patch -buN -p1 -d $RTC_DIR < $PATCH_DIR/mute.diff

SRC_DIR=$PATCH_DIR/../src
SDK_DIR=$RTC_DIR/sdk/objc

echo "cp $SRC_DIR/components/audio/* $SDK_DIR/components/audio"
cp $SRC_DIR/components/audio/* $SDK_DIR/components/audio

echo "cp $SRC_DIR/native/api/* $SDK_DIR/native/api"
cp $SRC_DIR/native/api/* $SDK_DIR/native/api

echo "cp $SRC_DIR/native/src/audio/* $SDK_DIR/native/src/audio"
cp $SRC_DIR/native/src/audio/* $SDK_DIR/native/src/audio
