# This file contains Docker development variables. it should be included by the Makefile
# one directory above
CONTAINER_REPO_PATH := github.com/arschles/k8stress
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.9.1
DEV_ENV_WORK_DIR := /go/src/github.com/arschles/k8stress
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -e CGO_ENABLED=0 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}
