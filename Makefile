GO ?= go

MODULE_PATH := $(shell awk '/^module / { print $$2; exit }' $(dir $(abspath $(lastword $(MAKEFILE_LIST))))go.mod)

APP_NAMES := $(notdir $(patsubst %/,%,$(dir $(wildcard cmd/*/main.go))))

RELEASE_VERSION := $(shell git describe --tags --exact-match 2>/dev/null)
DEV_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
VERSION ?= $(if $(RELEASE_VERSION),$(RELEASE_VERSION),$(DEV_VERSION))
COMMIT := $(shell git rev-parse --short HEAD)

NAME = $(patsubst %-$(VERSION),%,$(@F))
LDFLAGS ?= '-X main.name=$(NAME) \
		   -X main.version=$(VERSION) \
		   -X main.commit=$(COMMIT) \
		   -s -w'

SRC_PATH = ./cmd/$(NAME)
DEST_PATHS = $(addprefix bin/,$(addsuffix -$(VERSION),$(APP_NAMES)))

export CGO_ENABLED = 0

install: $(APP_NAMES)

build: $(DEST_PATHS)

$(APP_NAMES):
	$(GO) install -ldflags $(LDFLAGS) $(SRC_PATH)

$(DEST_PATHS):
	$(GO) build -ldflags $(LDFLAGS) -o $@ $(SRC_PATH)

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

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

.PHONY: install generate deps build tidy update fmt test vet
