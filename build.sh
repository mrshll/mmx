#!bin/sh
rm -r docs/*
mkdir -p docs
mkdir -p docs/links
cp -r links/* docs/links/

mkdir -p docs/img
cp -r data/img/* docs/img/

cd src && go run main.go && cd ../