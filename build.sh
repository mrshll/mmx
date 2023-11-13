#!/bin/bash

CONTENT_DIR=$1

if [ -z $CONTENT_DIR ]; then
  echo "error: content directory not provided"
else

  for f in src/plugins/*.sh; do
    bash "$f" || break # execute successfully or break
  done

  rm $CONTENT_DIR/site/*.html
  rm -r $CONTENT_DIR/site/links

  mkdir -p $CONTENT_DIR/site/links
  cp -r links/* $CONTENT_DIR/site/links/

  cp -r static/. $CONTENT_DIR/site/

  cd src && lua mmx.lua $CONTENT_DIR && cd ../
fi
