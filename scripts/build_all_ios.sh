#!/bin/sh

# usage: ./build_all_ios.sh CONFIG_DIR

NOFETCH=

while getopts ":-:" OPT
do
  case $OPT in
    -)
      case "${OPTARG}" in
        nofetch) NOFETCH=true;;
      esac
      ;;
  esac
done
shift $((OPTIND - 1))

DTOOLS_DIR=$(cd $(dirname $0)/../build/depot_tools && pwd)
CONFIG_DIR=$(cd $1 && pwd)
SCRIPT_DIR=$(dirname $0)

if [ "$NOFETCH" != "true" ]; then
  $SCRIPT_DIR/fetch_depot_tools.sh $DTOOLS_DIR
  $SCRIPT_DIR/fetch_webrtc.sh $CONFIG_DIR
fi
$SCRIPT_DIR/build_ios.sh $CONFIG_DIR

