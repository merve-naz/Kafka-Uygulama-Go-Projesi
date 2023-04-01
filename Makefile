# golang build, install, and clean Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=main

# docker parameters
DOCKER_IMAGE_PREFIX=bb/
DOCKER_IMAGE_NAME=certificate-service
# read VERSION file and assign to DOCKER_IMAGE_VERSION
DOCKER_IMAGE_VERSION=$(shell cat VERSION)
OUT_PLATFORM=linux/amd64

all: build
go-build:
	$(GOBUILD) -o $(BINARY_NAME) -v
go-test:
	$(GOCMD) test -v ./...
go-clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
go-run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
docker-build:
	docker buildx build --platform $(OUT_PLATFORM) -t $(DOCKER_IMAGE_PREFIX)$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION) .
docker-up:
	docker-compose up -d