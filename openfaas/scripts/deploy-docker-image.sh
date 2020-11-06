#!/bin/sh

set -e

# Script Arguments
DOCKER_IMAGE=${1}
DOCKER_USERNAME=${DOCKER_HUB_USERNAME:-$2}
DOCKER_PASSWORD=${DOCKER_HUB_PASSWORD:-$3}

# Login to docker registry
echo "${DOCKER_PASSWORD}" | docker login --username "${DOCKER_USERNAME}" --password-stdin

# Push the image
docker push "$DOCKER_IMAGE"
