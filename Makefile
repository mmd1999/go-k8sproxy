APP_NAME := go-k8sproxy
REPO_NAME := mmd1999
CHART_NAME := k8sproxy
VERSION=$(shell git rev-parse --short HEAD)

HELM_EXPERIMENTAL_OCI=1

ifndef CHART_VERSION
$(error CHART_VERSION is not set)
endif

.PHONY: build publish tag scan helm-version helm-package helm-push

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

helm-version:
	@echo "Updating Helm version in Chart.yaml"
	@./scripts/chart_version.sh $(CHART_VERSION) $(VERSION)

helm-package: helm-version
	@echo "Packaging $(CHART_NAME):$(CHART_VERSION)"
	helm package charts/$(CHART_NAME) --version $(CHART_VERSION) --app-version $(VERSION) .

helm-push: helm-package
	@echo "Pushing $(CHART_NAME):$(CHART_VERSION) to oci://ghcr.io/$(REPO_NAME)"
	helm push mychart-$(CHART_VERSION).tgz oci://ghcr.io/$(REPO_NAME)/$(CHART_NAME)