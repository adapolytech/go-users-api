IMAGE_NAME := users-service
IMAGE_TAG ?= latest

start:
	@echo "Building docker image $(IMAGE_NAME):$(IMAGE_TAG)"

build: start
	docker build -t $(IMAGE_TAG):$(IMAGE_TAG) .

tag-remote: build
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REPOSITORY_NAME)/$(IMAGE_NAME):$(IMAGE_TAG)

login:
	docker login -u $(DOCKERHUB_USERNAME) -p $(DOCKERHUB_PASSWORD) registry.hub.docker.com

push: login
	docker push $(REPOSITORY_NAME)/$(IMAGE_NAME):$(IMAGE_TAG)
