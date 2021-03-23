#!/bin/bash

rm -rf output
for arch in arm amd64; do
  for os in linux; do
    OUTDIR=output/$os/$arch
    echo "Building $os $arch"
    mkdir -p $OUTDIR
    export GOOS=$os
    export GOARCH=$arch
    go build -o $OUTDIR/goruuvitag github.com/joelmertanen/goruuvitag
  done
done

echo "Creating tar file"
tar zcvf binaries.tar.gz output
