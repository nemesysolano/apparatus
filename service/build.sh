#!/usr/bin/env bash

set -e
# https://www.youtube.com/watch?v=HelHcqbTzmA
if [ ! -d "bin" ]; then
  echo "Creating bin directory..."
  mkdir -p bin
else
  echo "Bin directory already exists, cleaning up..."
  rm  -rf bin/*
fi

go build -o bin/apparatus
STATUS=$?
if [ $STATUS -ne 0 ]; then
  echo "Build failed with status $STATUS"
  exit $STATUS
fi

