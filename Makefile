PROJECT=$(shell basename $(shell pwd))
TAG=ghcr.io/johnjones4/${PROJECT}
VERSION=$(shell date +%s)

.PHONY: ui

info:
	echo ${PROJECT} ${VERSION}

container:
	docker build --platform linux/amd64 -t ${TAG} ./server
	docker push ${TAG}:latest
	docker image rm ${TAG}:latest

ci: container
