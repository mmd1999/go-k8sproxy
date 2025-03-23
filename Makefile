APP_NAME := go-k8sproxy
REPO_NAME := mmd1999

VERSION=$(shell git rev-parse --short HEAD)

.PHONY: build publish tag scan

build:
	@echo "Building $(APP_NAME):$(VERSION)"
	docker build -t $(APP_NAME):$(VERSION) .

scan:
	@echo "Scanning $(REPO_NAME)/$(APP_NAME):$(VERSION)"
	trivy image --exit-code 1 -s "MEDIUM,HIGH,CRITICAL" $(APP_NAME):$(VERSION)

tag:
	@echo "Tagging $(REPO_NAME)/$(APP_NAME):$(VERSION)"
	docker tag $(APP_NAME):$(VERSION) $(REPO_NAME)/$(APP_NAME):$(VERSION)
	docker tag $(APP_NAME):$(VERSION) $(REPO_NAME)/$(APP_NAME):latest

publish: build scan tag
	@echo "Publishing $(REPO_NAME)/$(APP_NAME):$(VERSION)"
	docker push $(REPO_NAME)/$(APP_NAME):$(VERSION)
	docker push $(REPO_NAME)/$(APP_NAME):latest
