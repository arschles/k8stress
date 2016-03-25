# This file contains Docker production variables. It should be included by the Makefile
# one directory above.
DOCKER_REGISTRY ?= quay.io/
DOCKER_IMAGE_PREFIX ?= arschles
DOCKER_IMAGE_SHORT_NAME ?= k8stress
DOCKER_VERSION ?= git-$(shell git rev-parse --short HEAD)
DOCKER_IMAGE := ${DOCKER_REGISTRY}${DOCKER_IMAGE_PREFIX}/${DOCKER_IMAGE_SHORT_NAME}:${VERSION}
