#!/bin/sh
bash build.sh
while inotifywait -qqre modify ./src ./links ./data; do
  bash build.sh
done
