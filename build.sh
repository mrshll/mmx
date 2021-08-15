#!/bin/bash

for f in src/plugins/*.sh; do
  bash "$f" || break  # execute successfully or break
done

rm docs/*.html
rm -r docs/links

mkdir -p docs/links
cp -r links/* docs/links/

cp CNAME docs/CNAME

cd src && go run main.go mmxup.go && cd ../
