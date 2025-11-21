PROJECT=$(shell basename $(shell pwd))
TAG=ghcr.io/johnjones4/${PROJECT}
VERSION=$(shell date +%s)

.PHONY: ui

info:
	echo ${PROJECT} ${VERSION}

container:
	podman build -t ${TAG} ./server
	podman push ${TAG}:latest
	podman image rm ${TAG}:latest

ci: container
