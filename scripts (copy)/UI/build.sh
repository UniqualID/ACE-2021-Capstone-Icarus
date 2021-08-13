#!/bin/bash
echo "Building all Go projects..."
for f in ./*.go; do
    go build $f
done 
echo "Done"
