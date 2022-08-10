#!/usr/bin/env bash
set -e

if [ -z "$FROM_MAKEFILE" ]; then
    echo "Do not call this file directly - use the make command"
    exit 1
fi

CONTAINER_RUNTIME=docker # Forcing docker

if [ "$CONTAINER_RUNTIME" == "docker" ]; then
    $CONTAINER_RUNTIME buildx create --use --name bazaar_builder
fi

cp -r $(pwd)/dist/bazaar_linux_arm_7 $(pwd)/dist/bazaar_linux_armv7
cp -r $(pwd)/dist/bazaar_linux_amd64_v1 $(pwd)/dist/bazaar_linux_amd64

$CONTAINER_RUNTIME buildx build -f ${CONTAINERFILE_NAME} --platform=${BUILDX_PLATFORMS} --build-arg "ALPINE_VERSION=${CONTAINER_ALPINE_VERSION}" ${CONTAINER_BUILDX_OPTIONS} . | tee 2> /dev/null

# rm -rf $(pwd)/dist/bazaar_linux_armv7
# rm -rf $(pwd)/dist/bazaar_linux_amd64

if [ "$CONTAINER_RUNTIME" == "docker" ]; then
    $CONTAINER_RUNTIME buildx rm bazaar_builder
fi
