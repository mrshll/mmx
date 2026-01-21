#!/bin/bash
CONTENT_DIR=$1
SITE_DIR=$2

if [ -z $CONTENT_DIR ] || [ -z $SITE_DIR ]; then
  echo "Usage: ./run.sh path/to/content path/to/site"
  exit 1
fi

bash build.sh $CONTENT_DIR $SITE_DIR

if [[ $OSTYPE == 'darwin'* ]]; then
  while fswatch -1 --exclude '\.git' --exclude '\.obsidian' $CONTENT_DIR; do
    bash build.sh $CONTENT_DIR $SITE_DIR
  done
else
  while inotifywait -qqre modify --exclude '\.git' --exclude '\.obsidian' $CONTENT_DIR; do
    bash build.sh $CONTENT_DIR $SITE_DIR
  done
fi
