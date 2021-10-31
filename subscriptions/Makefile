GO_FILES=$(shell go list ./...)
BUILD_OUTPUT ?= "bin/subscriptions"

BUILD_DATE=$$(date +%Y-%m-%d-%H:%M)
BUILD_VERSION=0.0.1
BUILD_DATE_VAR := github.com/snooyen/animal-facts/subscriptions/pkg/version.BuildDate
BUILD_VERSION_VAR := github.com/snooyen/animal-facts/subscriptions/pkg/version.BuildVersion
LDFLAGS ?= "-X $(BUILD_DATE_VAR)=$(BUILD_DATE) -X $(BUILD_VERSION_VAR)=$(BUILD_VERSION)"

DOCKER_IMAGE = snooyen/subscriptions
DOCKER_CLI_EXPERIMENTAL = enabled
export GOPATH DOCKER_CLI_EXPERIMENTAL

.PHONY: build
build: setup
	@echo "Building for OSX..."
	$Q GOOS=darwin CGO_ENABLED=0 go build -ldflags ${LDFLAGS} -o ${BUILD_OUTPUT} ./cmd/

.PHONY: build-linux
build-linux: setup
	@echo "Building for LINUX..."
	$Q GOOS=linux CGO_ENABLED=0 go build -ldflags ${LDFLAGS} -o ${BUILD_OUTPUT} ./cmd/

image:
	@echo "Building docker image..."
	docker buildx build --no-cache -t $(DOCKER_IMAGE):$(BUILD_VERSION) --build-arg LDFLAGS=${LDFLAGS} -f Dockerfile --platform="linux/amd64" --output "type=docker,push=false" .

setup: clean
	@echo "Setup..."
	$Q go env

clean:
	-rm -rf bin

.PHONY: mod
mod:
	@go mod tidy -compat=1.17

.PHONY: unit-test
unit-test:
	go test $(GO_FILES) -v

.PHONY: coverage
coverage:
	go test $(GO_FILES) -cover

.PHONY: coverprofile
coverprofile:
	go test $(GO_FILES) -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: fmt
fmt:
	go fmt $(GO_FILES)
