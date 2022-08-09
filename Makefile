PROJECT := bazaar
VERSION_COMMIT := $git rev-parse --short HEAD)
SOURCE_FILES ?=./internal/... ./cmd/... ./pkg/...
TEST_OPTIONS ?= -v -failfast -race -bench=. -benchtime=100000x -cover -coverprofile=coverage.out
TEST_TIMEOUT ?=1m
CLEAN_OPTIONS ?=-modcache -testcache
LD_FLAGS := -X main.version=$(VERSION) -s -w
BUILDS_PATH := ./build
FROM_MAKEFILE := y

# Common exports
export FROM_MAKEFILE
export VERSION_COMMIT

export LD_FLAGS
export SOURCE_FILES
export TEST_OPTIONS
export TEST_TIMEOUT
export BUILDS_PATH

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
	@go build -a -v -ldflags "${LD_FLAGS}" -o ${BUILDS_PATH}/bazaar ./cmd/bazaar/*.go
	@chmod +x ${BUILDS_PATH}/bazaar

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
