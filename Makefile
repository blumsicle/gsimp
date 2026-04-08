GO ?= go

MODULE_PATH := $(shell awk '/^module / { print $$2; exit }' $(dir $(abspath $(lastword $(MAKEFILE_LIST))))go.mod)

APP_NAMES := $(notdir $(patsubst %/,%,$(dir $(wildcard cmd/*/main.go))))

RELEASE_VERSION := $(shell git describe --tags --exact-match 2>/dev/null)
DEV_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
VERSION ?= $(if $(RELEASE_VERSION),$(RELEASE_VERSION),$(DEV_VERSION))

NAME = $(patsubst %-$(VERSION),%,$(@F))

SRC_PATH = ./cmd/$(NAME)
DEST_PATHS = $(addprefix bin/,$(addsuffix -$(VERSION),$(APP_NAMES)))

export CGO_ENABLED = 0

install: $(APP_NAMES)

build: $(DEST_PATHS)

rebuild:
	$(MAKE) -B build

$(APP_NAMES):
	$(GO) install $(SRC_PATH)

$(DEST_PATHS):
	$(GO) build -o $@ $(SRC_PATH)

generate:
	$(GO) generate ./...

deps:
	$(GO) mod download

tidy:
	$(GO) mod tidy

update:
	$(GO) get -u ./...

fmt:
	gofumpt -w .
	goimports -w .
	golines -w .

check: vet test

clean:
	rm -rf bin

test:
	$(GO) test ./...

coverage:
	$(GO) test -cover ./...

coverage-html:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

vet:
	$(GO) vet ./...

.PHONY: install generate deps build rebuild tidy update fmt check clean test coverage coverage-html vet
