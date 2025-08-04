#!/usr/bin/env bash

set -e

if [ ! -d "bin" ]; then
  echo "Creating bin directory..."
  mkdir -p bin
else
  echo "Bin directory already exists, cleaning up..."
  rm  -rf bin/*
fi

mkdir bin/resources
cp -r resources/* bin/resources/

go build -o bin/downloader
STATUS=$?
if [ $STATUS -ne 0 ]; then
  echo "Build failed with status $STATUS"
  exit $STATUS
fi

