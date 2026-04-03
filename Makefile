MODULE_PATH := $(shell awk '/^module / { print $$2; exit }' $(dir $(abspath $(lastword $(MAKEFILE_LIST))))go.mod)

PROJ_NAME := $(notdir $(MODULE_PATH))
APP_NAMES := $(notdir $(dir $(wildcard cmd/*/*.go)))

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(shell basename $(BRANCH))
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

$(APP_NAMES): generate
	go install -ldflags $(LDFLAGS) $(SRC_PATH)

$(DEST_PATHS): generate
	go build -ldflags $(LDFLAGS) -o $@ $(SRC_PATH)

generate:
	go generate ./...

deps:
	go mod download

tidy:
	go mod tidy

update:
	go get -u ./...

fmt:
	gofumpt -w .

vet:
	go vet ./...

.PHONY: install generate deps build tidy update fmt vet
