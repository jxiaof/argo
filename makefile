default_target: run

DOCKER_DEPOS := 120.131.2.186:8080
IMAGE_NAME := iris-web

IMAGE_TAG := dev
IMAGE_PATH := test

IMAGE_FATHER_PATH := $(DOCKER_DEPOS)/$(IMAGE_PATH)

# 是否推送到远程
PUSH :=


.PHONY: docker
docker:
		docker build -f Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) .
ifdef PUSH
	    docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_FATHER_PATH)/$(IMAGE_NAME):$(IMAGE_TAG)
	    docker push $(IMAGE_FATHER_PATH)/$(IMAGE_NAME):$(IMAGE_TAG)
endif


.PHONY: docker.prod
docker.prod:
		$(MAKE) docker IMAGE_PATH=production


.PHONY: run
run:
		docker-compose -f docker-compose.yaml up
