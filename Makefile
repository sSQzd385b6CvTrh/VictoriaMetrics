# Makefile for VictoriaMetrics
# Provides common build, test, and deployment targets

APP_NAME := victoria-metrics
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GO ?= go
GOFLAGS ?= -trimpath
LDFLAGS := -X 'github.com/VictoriaMetrics/VictoriaMetrics/lib/buildinfo.Version=$(VERSION)' \
           -X 'github.com/VictoriaMetrics/VictoriaMetrics/lib/buildinfo.BuildTime=$(BUILD_TIME)'

# Personal fork: using a local registry prefix for docker images
DOCKER_IMAGE ?= localhost:5000/victoria-metrics
DOCKER_TAG ?= $(VERSION)

.PHONY: all
all: build

## build: Compile the application binary
.PHONY: build
build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) ./app/victoria-metrics

## build-all: Build all application binaries
.PHONY: build-all
build-all:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" ./app/...

## test: Run unit tests
.PHONY: test
test:
	$(GO) test ./... -count=1 -race -timeout 120s

## test-short: Run short unit tests
.PHONY: test-short
test-short:
	$(GO) test ./... -count=1 -short -timeout 60s

## bench: Run benchmarks
.PHONY: bench
bench:
	$(GO) test ./... -bench=. -benchmem -run='^$' -count=1

## lint: Run golangci-lint
.PHONY: lint
lint:
	golangci-lint run ./...

## fmt: Format Go source files
.PHONY: fmt
fmt:
	gofmt -w -s ./
	$(GO) mod tidy

## vet: Run go vet
.PHONY: vet
vet:
	$(GO) vet ./...

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		.

## docker-push: Push Docker image to registry
.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

## clean: Remove build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	$(GO) clean -cache

## mod-download: Download Go module dependencies
.PHONY: mod-download
mod-download:
	$(GO) mod download

## mod-verify: Verify Go module dependencies
.PHONY: mod-verify
mod-verify:
	$(GO) mod verify

## help: Display this help message
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
