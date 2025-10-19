#!/usr/bin/env bash

# Compile Go projects to WebAssembly (wasm) in all subdirectories that contain Go source files.
for dir in */ ; do
    echo "Building $dir wasm"
    (cd "$dir/wasm" && GOOS=js GOARCH=wasm go build -ldflags "-X main.Token=$MAPBOX_ACCESS_TOKEN" -o main.wasm)
done