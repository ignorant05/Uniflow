.PHONY: install build test ctest fmt lint clean build-docker buildx-setup login-ghcr ghcr push help  

ARGS = main.go
OUT = uniflow 
TEST_ARGS = ./...

IMAGE_NAME = uniflow
GITHUB_USER = ignorant05
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "latest")

TAG = ghcr.io/$(GITHUB_USER)/$(IMAGE_NAME):$(VERSION)
LATEST = ghcr.io/$(GITHUB_USER)/$(IMAGE_NAME):latest

install: 
	@echo "Installing dependencies..."
	go mod tidy 
	@echo "Dependencies Installed Successfully"
	@echo "You can proceed now..."

test: 
	@echo "Running tests..."
	go test $(TEST_ARGS) -v
	@echo "Tests ran"

ctest: 
	@echo "Running coverage tests..."
	go test -cover $(TEST_ARGS) -v
	@echo "Tests ran"

build: test 
	@echo "Building Uniflow..."
	go build -ldflags="-w -s" -trimpath -o $(OUT) .
	@echo "Uniflow Built Successfully"

fmt:
	@echo "Running go fmt..."
	go fmt $(TEST_ARGS) 
	@echo "Formatting completed"

lint:
	@echo "Running go vet..."
	go vet $(TEST_ARGS)
	@echo "Running staticcheck"
	@which staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "staticcheck not installed, skipping"
	@echo "Linting completed"

clean: 
	@echo "Cleaning Up..."
	rm -f $(OUT) 
	@echo "Old binary cleaned up"
	docker rmi $(TAG) 2>/dev/null || true
	@echo "Docker image removed"

build-docker: 
	@echo "Building image TAG: $(TAG) ..."
	docker build -t $(TAG) .
	@echo "Image built"

buildx-setup:
	@echo "Setting up docker buildx ..."
	docker buildx create --name mybuilder --use
	docker buildx inspect --bootstrap
	@echo "Docker buildx sat up"

login-ghcr: 
	@if [ -z "$(GHCR_TOKEN)" ]; then \
		echo "Error: GHCR_TOKEN env var is not sat"; \
		exit 1; \
	fi  
	echo "$(GHCR_TOKEN)" | docker login ghcr.io -u $(GITHUB_USER) --password-stdin

push: login-ghcr buildx-setup
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--tag $(TAG) \
		--tag $(LATEST) \
		--push .

ghcr: test push
	@echo "Pushed to GHCR:"
	@echo "   $(TAG)"
	@echo "   $(LATEST)"

help:
	@echo "Available commands:"
	@echo "  install      - Resolve dependencies"
	@echo "  test         - Run tests"
	@echo "  ctest        - Run coverage tests"
	@echo "  build        - Build local binary"
	@echo "  fmt          - Code formatting"
	@echo "  lint         - Code linting"
	@echo "  clean        - Clean build artifacts"
	@echo "  build-docker - Build Docker image locally"
	@echo "  push         - Build and push multi-arch image to GHCR"
	@echo "  ghcr         - Test, build and push to GHCR"
	@echo ""
	@echo "Usage:"
	@echo "  make ghcr GHCR_TOKEN=ghp-token" 
	@echo "  make push VERSION=0.2.0"
