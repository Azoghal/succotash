#!/bin/bash

set -e -o pipefail

cd ../

rm -rf build/webpage
mkdir -p build/webpage

cp -r frontend/dist build/webpage