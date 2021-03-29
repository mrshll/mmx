#!bin/sh
rm -r docs/*
mkdir -p docs
mkdir -p docs/links
cp -r links/* docs/links/

mkdir -p docs/img
cp -r data/img/* docs/img/

cp CNAME docs/CNAME

cd src && go run main.go mmxup.go && cd ../
