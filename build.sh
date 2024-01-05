#!/bin/bash

CONTENT_DIR=$1
SITE_DIR=$2

if [ -z $CONTENT_DIR ] || [ -z $SITE_DIR ]; then
  echo "Usage: ./build.sh path/to/content path/to/site"
  exit 1
fi

# TODO: skipping plugins for now
# for f in src/plugins/*.sh; do
#   bash "$f" $SITE_DIR || break # execute successfully or break
# done

rm $SITE_DIR/*.html 2>/dev/null
rm -r $SITE_DIR/links 2>/dev/null

mkdir -p $SITE_DIR/links
cp -r links/* $SITE_DIR/links/

cp -r static/. $SITE_DIR

cd src && lua mmx.lua $CONTENT_DIR $SITE_DIR && cd ../

# temporary to support neon kiosk
cp $SITE_DIR/Log.html $SITE_DIR/marshall_s_journal.html
