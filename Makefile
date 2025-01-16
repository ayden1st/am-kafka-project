DOCKER_USERNAME ?= ayden1st
APPLICATION_NAME ?= am-kafka-project
VERSION ?= 0.7-dev

GIT_COMMIT := $(shell git rev-list -1 HEAD)
GIT_TAG_VERSION := $(shell git tag --points-at HEAD)
BUILD_DATE := $(shell date -u +%Y%m%d-%H:%M:%S)

build:
	docker build --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg GIT_TAG_VERSION=$(GIT_TAG_VERSION) --build-arg BUILD_DATE=$(BUILD_DATE) --tag ${DOCKER_USERNAME}/${APPLICATION_NAME}:${VERSION} .

push:
	docker push ${DOCKER_USERNAME}/${APPLICATION_NAME}:${VERSION}
