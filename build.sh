#!bin/sh
mkdir -p docs
mkdir -p links
cp -r links/* docs/links/

mkdir -p img
cp -r data/img/* docs/img/

cd src && go run main.go && cd ../