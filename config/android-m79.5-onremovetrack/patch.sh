#!/bin/sh

patch -buN -p1 -d $RTC_DIR < $PATCH_DIR/android.diff

