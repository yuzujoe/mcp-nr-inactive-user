# Variables
APP_NAME = mcp-server
DOCKER_IMAGE = mcp-nr-inactive-user
DIST_DIR = dist
GO_FILES = $(shell find . -name '*.go' -not -path "./vendor/*")
TAG ?= latest

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o $(DIST_DIR)/$(APP_NAME) .

# Run the application locally
.PHONY: run
run: build
	./$(DIST_DIR)/$(APP_NAME)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(DIST_DIR)

# Docker targets
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(TAG) .

.PHONY: docker-run
docker-run: docker-build
	docker run --rm $(DOCKER_IMAGE)

# Format code
.PHONY: fmt
fmt:
	gofmt -w $(GO_FILES)

# Lint code
.PHONY: lint
lint:
	go vet ./...
	golint ./...
