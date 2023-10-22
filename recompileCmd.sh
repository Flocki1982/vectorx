#!/bin/bash

echo "Compiling VectorX to speed up execution"
echo
/usr/local/go/bin/go build cmd/main.go
mv main vectorx
