PROJECT_NAME := bazaar
GOLANG_VERSION := $(shell head -n 3 go.mod | tail -n 1 | cut -d " " -f 2)

SOURCE_FILES ?=./internal/... ./cmd/... ./pkg/...

TEST_OPTIONS ?= -v -failfast -race -bench=. -benchtime=100000x -cover -coverprofile=coverage.out
TEST_TIMEOUT ?=1m

CLEAN_OPTIONS ?=-modcache -testcache

CGO_ENABLED := 0

BUILDS_PATH := ./dist
FROM_MAKEFILE := y

CONTAINER_RUNTIME := podman
CONTAINERFILE_NAME := Containerfile
CONTAINER_ALPINE_VERSION := 3.19

BUILDX_PLATFORMS := linux/amd64,arm64,linux/arm/v7

# Common exports
export GOLANG_VERSION

export FROM_MAKEFILE

export CGO_ENABLED

export SOURCE_FILES
export TEST_OPTIONS
export TEST_TIMEOUT
export BUILDS_PATH

export CONTAINER_RUNTIME
export CONTAINERFILE_NAME
export CONTAINER_ALPINE_VERSION

export BUILDX_PLATFORMS

.PHONY: all
all: help

# this is godly
# https://news.ycombinator.com/item?id=11939200
.PHONY: help
help:	### this screen. Keep it first target to be default
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# this is not tested, but prepared in advance for you, Mac drivers
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

.PHONY: clean
clean: ###  clean test cache, build files
	$(info: Make: Clean)
	@rm -rf ${BUILDS_PATH}
	@go clean ${CLEAN_OPTIONS}

.PHONY: build
build: clean ### builds the project for the setup os/arch combinations
	$(info: Make: Build)
	@goreleaser build --clean --snapshot

# Only used in CI, assumes the build process has already been run
.PHONY: buildx
buildx:
	$(info: Make: Buildx)
	@bash scripts/buildx.sh

.PHONY: quick-run
quick-run: ### Executes the project using golang
	@go run ./cmd/bazaar/*.go

.PHONY: run
run: ### Executes the project build locally
	@make build
	${BUILDS_PATH}/bazaar

.PHONY: format
format: ### Executes the formatting pipeline on the project
	$(info: Make: Format)
	@bash scripts/format.sh

.PHONY: lint
lint: ### Check the project for errors
	$(info: Make: Lint)
	@bash scripts/lint.sh

.PHONY: test
test: ### Runs the test suite
	@bash scripts/test.sh
