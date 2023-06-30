DOCKER_USERNAME := pehli
APPLICATION_NAME := node-wizard

.PHONY: build
build:
	docker build --tag $(DOCKER_USERNAME)/${APPLICATION_NAME}:latest .

.PHONY: push
push:
	docker push $(DOCKER_USERNAME)/${APPLICATION_NAME}:latest
