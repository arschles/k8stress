include make/docker-dev.mk
include make/docker.mk

build:
	${DEV_ENV_CMD} go build -o rootfs/bin/k8stress

docker-build:
	docker build --rm -t ${DOCKER_IMAGE} rootfs

docker-push:
	docker push ${DOCKER_IMAGE}
