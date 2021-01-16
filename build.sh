#!bin/sh
mkdir -p doc
mkdir -p links
cp -r links/* doc/links/

mkdir -p img
cp -r data/img/* doc/img/

cd src && go run main.go && cd ../