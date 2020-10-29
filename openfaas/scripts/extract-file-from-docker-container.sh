#!/bin/bash

set -e

##############################################################################################
#
# Extract a file from a docker container onto the local file system
#
##############################################################################################

# Script Arguments
DOCKER_IMAGE=${1}
DOCKER_PATH=${2}
LOCAL_FILE_SYSTEM_PATH=${3}

# Create the local file system path if necessary
mkdir -p "$LOCAL_FILE_SYSTEM_PATH"

# Create a temporary container to copy from
CONTAINER_ID=$(docker create "$DOCKER_IMAGE")

# Copy the native images out of the docker container
docker cp "$CONTAINER_ID:${DOCKER_PATH}" "$LOCAL_FILE_SYSTEM_PATH"

# Remove the temporary container
docker rm -f -v "$CONTAINER_ID"
