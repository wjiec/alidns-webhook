IMAGE_NAME ?= wjiec/alidns-webhook
IMAGE_TAG ?= $(shell cat VERSION)

.PHONY: unit-test
unit-test:
	go test -v ./...

.PHONY: e2e-test
e2e-test:
	echo ok

.PHONY: build
build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) --push .
