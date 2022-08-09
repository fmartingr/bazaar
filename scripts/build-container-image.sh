#!/usr/bin/env bash
set -e

if [ -z "$FROM_MAKEFILE" ]; then
    echo "Do not call this file directly - use the make command"
    exit 1
fi

${CONTAINER_RUNTIME} build --build-arg "GOLANG_VERSION=${CONTAINER_GOLANG_VERSION}" --build-arg="ALPINE_VERSION=${CONTAINER_ALPINE_VERSION}" -t ${CONTAINER_IMAGE_NAME}:${CONTAINER_IMAGE_TAG} -f ${CONTAINERFILE_NAME} .
