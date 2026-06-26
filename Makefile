BINARY_NAME=supervisor
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -ldflags "\
	-X github.com/ram291/opamp-control-pane/internal/version.Version=$(VERSION) \
	-X github.com/ram291/opamp-control-pane/internal/version.CommitHash=$(COMMIT_HASH) \
	-X github.com/ram291/opamp-control-pane/internal/version.BuildDate=$(BUILD_DATE)"

.PHONY: all
all: clean deps ui-build build

.PHONY: build
build: ui-build
	@echo "Building supervisor binary..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/supervisor/
	@echo "Build complete: bin/$(BINARY_NAME)"

.PHONY: build-linux
build-linux: ui-build
	@echo "Building supervisor binary for linux/amd64..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/supervisor/
	@echo "Build complete: bin/$(BINARY_NAME)-linux-amd64"

.PHONY: clean
clean:
	rm -rf bin/
	rm -rf ui/dist/
	@echo "Clean complete"

.PHONY: test
test:
	go test -v -race -count=1 ./...
	@echo "Tests complete"

.PHONY: test-cover
test-cover:
	go test -v -race -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: lint
lint:
	golangci-lint run ./...
	@echo "Lint complete"

.PHONY: deps
deps:
	@echo "Downloading Go dependencies..."
	go mod download
	go mod tidy
	@echo "Go dependencies ready"

.PHONY: ui-build
ui-build:
	@echo "Building React frontend..."
	cd ui && npm ci --silent && npm run build
	@echo "React build complete"

.PHONY: ui-dev
ui-dev:
	@echo "Starting React dev server..."
	cd ui && npm run dev

.PHONY: dev
dev:
	@echo "Starting Go backend..."
	go run $(LDFLAGS) ./cmd/supervisor/

.PHONY: run
run: build
	./bin/$(BINARY_NAME)

.PHONY: docker
docker:
	docker build -t opamp-control-pane:$(VERSION) .
	@echo "Docker image built: opamp-control-pane:$(VERSION)"

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -v $(PWD)/configs:/etc/opamp opamp-control-pane:$(VERSION)

.PHONY: update-opamp
update-opamp:
	go get github.com/open-telemetry/opamp-go@latest
	go mod tidy
	@echo "Updated opamp-go to latest version"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Clean, install deps, build UI, build binary"
	@echo "  build         - Build the supervisor binary (with UI embedded)"
	@echo "  build-linux   - Cross-compile for linux/amd64"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run all tests"
	@echo "  test-cover    - Run tests with coverage report"
	@echo "  lint          - Run golangci-lint"
	@echo "  deps          - Download Go dependencies"
	@echo "  ui-build      - Build the React frontend"
	@echo "  ui-dev        - Start React dev server with HMR"
	@echo "  dev           - Run Go backend directly (without building UI)"
	@echo "  run           - Build and run the binary"
	@echo "  docker        - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  update-opamp  - Update opamp-go dependency to latest"