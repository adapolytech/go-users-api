IMAGE_NAME := users-service
IMAGE_TAG ?= latest

start:
	@echo "Building docker image $(IMAGE_NAME):$(IMAGE_TAG)"

build: start
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

tag-remote:
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REPOSITORY_NAME):$(IMAGE_TAG)

login:
	@echo $(DOCKERHUB_PASSWORD) | docker login -u $(DOCKERHUB_USERNAME) --password-stdin

push: login
	docker push $(REPOSITORY_NAME):$(IMAGE_TAG)
