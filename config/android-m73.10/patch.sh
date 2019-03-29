#!/bin/sh

patch -buN $RTC_DIR/tools_webrtc/android/build_aar.py $PATCH_DIR/build_aar.py.diff
patch -buN $RTC_DIR/sdk/android/BUILD.gn $PATCH_DIR/BUILD.gn.diff

