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

## check: Run fmt, vet, and lint together (useful before committing)
.PHONY: check
check: fmt vet lint

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

## test-verbose: Run unit tests with verbose output (handy for debugging)
.PHONY: test-verbose
test-verbose:
	$(GO) test ./... -count=1 -race -timeout 120s -v

## test-cover: Run tests with coverage report (personal addition for tracking coverage)
# Note: open coverage.html in a browser to explore uncovered lines interactively
.PHONY: test-cover
test-cover:
	$(GO) test ./... -count=1 -race -timeout 120s -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to coverage.html"
	@$(GO) tool cover -func=coverage.out | tail -1

## cover-open: Generate coverage report and open it in the default browser
# Personal convenience target — detects OS to pick the right open command
.PHONY: cover-open
cover-open: test-cover
	@if [ "$(shell uname)" = "Darwin" ]; then \
		open coverage.html; \
	elif [ "$(shell uname)" = "Linux" ]; then \
		xdg-open coverage.html; \
	else \
		echo "Open coverage.html manually in your browser"; \
	fi
