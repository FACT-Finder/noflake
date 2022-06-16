#!/usr/bin/env sh
DIR=$(dirname "$0")
DIST=$DIR/dist
if [ -d "$DIST" ]; then
    echo "$DIST exists, skipping swagger download"
    exit 0
fi

VERSION=4.12.0
rm -rf "$DIST"
mkdir -p "$DIST"
wget -O $DIR/swagger-ui.tar.gz https://github.com/swagger-api/swagger-ui/archive/refs/tags/v$VERSION.tar.gz
tar -xzf $DIR/swagger-ui.tar.gz --strip-components 2 -C "$DIST" swagger-ui-$VERSION/dist
rm swagger-ui.tar.gz
