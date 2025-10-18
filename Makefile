# Makefile para build com informações de versão
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse HEAD)
GO_VERSION := $(shell go version | awk '{print $$3}')

# Flags de build para injetar versão
LDFLAGS := -ldflags "\
	-X rest-api-2/internal/version.Version=$(VERSION) \
	-X rest-api-2/internal/version.BuildDate=$(BUILD_DATE) \
	-X rest-api-2/internal/version.GitCommit=$(GIT_COMMIT) \
	-X rest-api-2/internal/version.GoVersion=$(GO_VERSION)"

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/api cmd/api/main.go

.PHONY: build-client
build-client:
	go build $(LDFLAGS) -o bin/client cmd/client/main.go

.PHONY: run
run:
	go run $(LDFLAGS) cmd/api/main.go

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GO_VERSION)"

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: test
test:
	go test -v ./...

.PHONY: docker-build
docker-build:
	docker build -t rest-api-2:$(VERSION) .

.DEFAULT_GOAL := build