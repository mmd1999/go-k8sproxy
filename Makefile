APP_NAME := go-k8sproxy
REPO_NAME := mmd1999

VERSION=$(shell git rev-parse --short HEAD)

.PHONY: build publish

build:
	@echo "Building $(APP_NAME):$(VERSION)"
	docker build -t $(APP_NAME):$(VERSION) .

publish: build
	docker tag $(APP_NAME):$(VERSION) $(REPO_NAME)/$(APP_NAME):$(VERSION)
	docker login -u $(REPO_NAME)
	docker push $(REPO_NAME)/$(APP_NAME):$(VERSION)

tag:
	docker tag $(APP_NAME):$(VERSION) $(REPO_NAME)/$(APP_NAME):$(VERSION)
	git rev-parse --short HEAD

version:
	@echo $(VERSION)