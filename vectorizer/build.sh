#!/usr/bin/env bash

set -e
# https://www.youtube.com/watch?v=HelHcqbTzmA
if [ ! -d "bin" ]; then
  echo "Creating bin directory..."
  mkdir -p bin/resources
else
  echo "Bin directory already exists, cleaning up..."
  rm  -rf bin/*
fi

if [ -d "resources" ] && [ "$(ls -A resources)" ]; then
  cp -r resources/* bin/resources
fi

go build -o bin/vectorizer
STATUS=$?
if [ $STATUS -ne 0 ]; then
  echo "Build failed with status $STATUS"
  exit $STATUS
fi

