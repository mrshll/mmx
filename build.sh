#!bin/sh
mkdir -p site
mkdir -p links
cp -r links/* site/links/

mkdir -p img
cp -r data/img/* site/img/

cd src && go run main.go && cd ../