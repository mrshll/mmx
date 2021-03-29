#!/bin/sh
cd src
go test
while inotifywait -qqre modify . ; do
  go test
done
