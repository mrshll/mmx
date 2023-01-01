#!/bin/sh
bash build.sh
if [[ $OSTYPE == 'darwin'* ]]; then
  while fswatch -1 ./src ./links ./data; do
    bash build.sh
  done
else
  while inotifywait -qqre modify ./src ./links ./data; do
    bash build.sh
  done
fi
