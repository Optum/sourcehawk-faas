#!/bin/bash

set -e

##################################################################################
#
# Docker Build
#
# This script is only used for local development purposes
#
##################################################################################
ROOT="$(dirname "$(dirname "$(dirname "$(readlink -fm "$0")")")")"
TAG="optumopensource/sourcehawk:faas-scan-SNAPSHOT"
NATIVE_IMAGE_PATH="$ROOT/scan/target/native-image"

# Make sure the native image exists first
if [[ ! -f "$NATIVE_IMAGE_PATH" ]]; then
  echo "$NATIVE_IMAGE_PATH does not exist"
  exit 1
fi

# Build the docker image
echo "Building docker image as: $TAG"
docker build -t "$TAG" "$ROOT/scan"