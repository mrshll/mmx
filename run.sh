#!/bin/sh
cd src && go run main.go && cd ../
while inotifywait -qqre modify ./src ./links ./data; do
   cd src && go run main.go && cd ../
done