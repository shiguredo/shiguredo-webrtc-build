#!/bin/bash

# usage: generate_version_android.sh CONFIG_DIR

CONFIG_DIR=$(cd $1 && pwd)
BUILD_CONFIG_FILE=$CONFIG_DIR/CONFIG
VERSION_CONFIG_FILE=$CONFIG_DIR/VERSION

_BUILD_DIR=$(dirname $0)/../build/$(basename $CONFIG_DIR)
mkdir -p $_BUILD_DIR
BUILD_DIR=$(cd $_BUILD_DIR && pwd)

SCRIPT_DIR=$(cd $(dirname $0) && pwd)
RTC_DIR=$BUILD_DIR/src

export PATH=$DTOOLS_DIR:$PATH

source $VERSION_CONFIG_FILE

cat <<EOF
package org.webrtc;

public interface WebrtcBuildVersion {
    public static final String webrtc_branch = "M$BRANCH";
    public static final String webrtc_commit = "$COMMIT";
    public static final String webrtc_revision = "$REVISION";
    public static final String maint_version = "$MAINT";
}
EOF

