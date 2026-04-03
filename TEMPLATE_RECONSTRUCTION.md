# Template Reconstruction Guide

This document describes the current repository in enough detail to recreate a very similar project from scratch in a new directory, even in a fresh Codex session.

It is intentionally more explicit than the main `README.md`. The goal is not just to explain how to use the repo, but to preserve the structure, rationale, and exact wiring choices that define the current template.

## Purpose

This repository is a starter template for Go CLI tools.

Core goals:

- provide a thin `main` for each binary
- centralize CLI parsing and logger wiring in a shared package
- keep command-specific flags and command handlers easy to edit
- support linker-injected build metadata
- support multiple binaries in one repo if needed
- provide a Makefile with common local workflows
- include smoke tests for the wiring most likely to regress

Current stack:

- CLI parser: `github.com/alecthomas/kong`
- logging: `github.com/rs/zerolog`
- tests: `github.com/stretchr/testify`

## Current Files

Current tracked source files:

- `go.mod`
- `go.sum`
- `Makefile`
- `README.md`
- `TEMPLATE_RECONSTRUCTION.md`
- `cmd/globals.go`
- `cmd/gsimp/main.go`
- `cmd/gsimp/cli.go`
- `cmd/gsimp/example/cmd.go`
- `cmd/gsimp/main_test.go`
- `internal/cli/buildinfo.go`
- `internal/cli/runner.go`

## Module and Dependencies

`go.mod` currently declares:

- module path: `github.com/blumsicle/gsimp`
- Go version: `1.26.1`

Direct dependencies:

- `github.com/alecthomas/kong v1.15.0`
- `github.com/rs/zerolog v1.35.0`
- `github.com/stretchr/testify v1.11.1`

Indirect dependencies currently present:

- `github.com/davecgh/go-spew v1.1.1`
- `github.com/mattn/go-colorable v0.1.14`
- `github.com/mattn/go-isatty v0.0.20`
- `github.com/pmezard/go-difflib v1.0.0`
- `golang.org/x/sys v0.42.0`
- `gopkg.in/yaml.v3 v3.0.1`

## Package Layout

The package layout is deliberate:

- `cmd/<binary>` is for binary entrypoints and the CLI tree for that binary.
- `cmd/globals.go` defines shared injected arguments used by command handlers.
- `internal/cli` contains shared runtime infrastructure for parsing, build metadata, and logger creation.

Why this split exists:

- `internal/cli` is reusable framework-style code.
- `cmd/globals.go` is app-level injected data, not parser/runtime infrastructure.
- command packages such as `cmd/gsimp/example` can import `github.com/blumsicle/gsimp/cmd` and receive the same concrete `*cmd.Globals` type during Kong handler injection.

This is important: the project intentionally does not keep `Globals` in `main`, because subcommands cannot import package `main`, and Kong injection requires the exact same concrete type on both sides.

## CLI Architecture

The current CLI architecture has three layers.

### 1. Shared runtime layer

Files:

- `internal/cli/buildinfo.go`
- `internal/cli/runner.go`

Responsibilities:

- define build metadata shape
- construct Kong options from config
- construct a Kong parser
- parse `os.Args`
- construct a zerolog logger
- execute the selected Kong command via `ctx.Run(...)`

Types and functions:

- `type BuildInfo struct { Name, Version, Commit string }`
- `type Config struct { Description string; BuildInfo BuildInfo }`
- `type Runner interface { GetLogLevel() zerolog.Level; RunArgs() []any }`
- `func Options(cfg Config) []kong.Option`
- `func New(app any, cfg Config, options ...kong.Option) (*kong.Kong, error)`
- `func Parse(app any, cfg Config) *kong.Context`
- `func NewLogger(level zerolog.Level) zerolog.Logger`
- `func Run(ctx *kong.Context, log zerolog.Logger, args ...any) error`

Important behavior:

- `Options` sets:
  - Kong app name from `cfg.BuildInfo.Name`
  - Kong description from `cfg.Description`
  - compact help formatting
  - the `version` interpolation variable
- `Parse` uses `os.Args[1:]`
- `Parse` calls `parser.FatalIfErrorf(err)` if parse fails
- `NewLogger` returns a console logger with:
  - output to `os.Stderr`
  - console time format `time.DateTime + " MST"`
  - level set by `.Level(level)`
  - timestamp
  - a `logger=main` field
- `Run` prepends the logger to the injected args and calls `ctx.Run(...)`

### 2. Shared injected args layer

File:

- `cmd/globals.go`

Current contents:

- `type Globals struct { ConfigFile string ... }`

The current flag is:

- field name: `ConfigFile`
- short flag: `-c`
- long flag: `--config-file`
- default: `~/.config/starter/config.yaml`
- Kong type: `path`
- help text: `Path to the config file`

This file is meant to be easy to customize when cloning the template.

### 3. Binary-specific CLI tree

Files:

- `cmd/gsimp/main.go`
- `cmd/gsimp/cli.go`
- `cmd/gsimp/example/cmd.go`

#### `cmd/gsimp/main.go`

Responsibilities:

- define linker-overridable build variables:
  - `name`
  - `version`
  - `commit`
- define the CLI description
- call the shared parser/runtime
- set zerolog global formatting knobs that are process-level:
  - `zerolog.DurationFieldUnit = time.Minute`
  - `zerolog.TimeFieldFormat = time.DateTime + " MST"`
- create the logger once
- log and exit on runtime error

Important design choice:

- `main` is intentionally thin
- process-level side effects stay in `main`
- shared parsing/runtime remains in `internal/cli`

Current description:

- `Starter CLI template`

#### `cmd/gsimp/cli.go`

Defines the root CLI struct.

Current structure:

- embeds `cmd.Globals`
- defines:
  - `LogLevel zerolog.Level`
  - `Version kong.VersionFlag`
  - `Example example.Command`

Current help texts:

- log level: `Log level`
- version flag: `Output version`
- example subcommand: `Example subcommand for new projects`

Interface methods implemented:

- `GetLogLevel() zerolog.Level`
- `RunArgs() []any`

`RunArgs()` currently returns:

- `[]any{&c.Globals}`

That shape matters because Kong handler injection uses these runtime-provided values when it invokes command `Run(...)` methods.

#### `cmd/gsimp/example/cmd.go`

Placeholder subcommand.

Current behavior:

- package name: `example`
- type: `Command`
- method: `func (c *Command) Run(log zerolog.Logger, g *cmd.Globals) error`
- logs:
  - message: `example command`
  - field: `config_file=<resolved config path>`

This command exists primarily as a starter example and as a test target for Kong injection.

## Build Metadata and Linker Flags

Build metadata is injected into `main` package variables via `-ldflags -X`.

Variables defined in `cmd/gsimp/main.go`:

- `name = "gsimp"`
- `version = "dev"`
- `commit = "unknown"`

The Makefile overwrites them with:

- `-X main.name=$(NAME)`
- `-X main.version=$(VERSION)`
- `-X main.commit=$(COMMIT)`

Important detail:

- for this project layout, `-X main.<var>` is the correct target
- earlier attempts to target import paths like `github.com/.../cmd/gsimp.name` did not work for the built binary

## Makefile Behavior

The Makefile is part of the template design, not an afterthought.

### Variables

- `GO ?= go`
- `MODULE_PATH` is extracted from `go.mod` using `awk`
- `APP_NAMES` is derived from `cmd/*/main.go`
- `RELEASE_VERSION` uses `git describe --tags --exact-match`
- `DEV_VERSION` uses `git describe --tags --always --dirty`
- `VERSION ?= $(if $(RELEASE_VERSION),$(RELEASE_VERSION),$(DEV_VERSION))`
- `COMMIT := $(shell git rev-parse --short HEAD)`
- `NAME = $(patsubst %-$(VERSION),%,$(@F))`
- `SRC_PATH = ./cmd/$(NAME)`
- `DEST_PATHS = $(addprefix bin/,$(addsuffix -$(VERSION),$(APP_NAMES)))`

### Version resolution rules

The project intentionally distinguishes release and local/dev versioning:

- if `HEAD` is exactly on a tag, `VERSION` becomes that tag
- otherwise `VERSION` becomes a Git description like `abc1234-dirty`
- the user or CI may still override with `VERSION=...`

### Targets

- `install`
  - installs binaries with `go install`
  - target list comes from `$(APP_NAMES)`
- `build`
  - incremental build into versioned files under `bin/`
  - uses `$(DEST_PATHS)` file targets
- `rebuild`
  - force rebuilds by running `$(MAKE) -B build`
- `generate`
  - `go generate ./...`
- `deps`
  - `go mod download`
- `tidy`
  - `go mod tidy`
- `update`
  - `go get -u ./...`
- `fmt`
  - `gofumpt -w .`
- `test`
  - `go test ./...`
- `vet`
  - `go vet ./...`
- `check`
  - runs `test` and `vet`
- `clean`
  - `rm -rf bin`

Important behavior:

- `build` is phony, but its work is delegated to file targets in `$(DEST_PATHS)`, so it can still be a no-op if those artifacts are up to date
- `rebuild` exists specifically to force rebuilding the versioned artifacts

### Binary naming

The build artifacts are versioned:

- output path pattern: `bin/<app>-<version>`

Examples:

- `bin/gsimp-v1.2.3`
- `bin/gsimp-a8b21fc-dirty`

## Testing Strategy

There is currently one test file:

- `cmd/gsimp/main_test.go`

It uses:

- `testing`
- `testify/assert`
- `testify/require`

### What is tested

The tests are smoke tests for the template wiring.

1. `TestVersionFlag`

- constructs a parser with shared runtime code
- supplies custom Kong writers
- overrides Kong `Exit` to capture the exit code
- parses `--version`
- asserts:
  - an error is returned
  - exit code is `0`
  - stdout contains `gsimp test-version test-commit`
  - stderr is empty

Important note:

- the test expects an error for `--version`
- this is intentional and matches Kong’s behavior when `Exit` is overridden in tests on a CLI with required subcommands

2. `TestHelpFlag`

- same parser setup as above
- parses `--help`
- asserts:
  - an error is returned
  - exit code is `0`
  - stdout includes:
    - `Starter CLI template`
    - `--config-file`
    - `example`
  - stderr is empty

3. `TestExampleCommandReceivesInjectedGlobals`

- parses:
  - `--config-file /tmp/test-config.yaml example`
- runs the selected command via `cliutil.Run(...)`
- uses a logger writing to a bytes buffer
- asserts:
  - no exit occurred
  - logs contain:
    - `example command`
    - `/tmp/test-config.yaml`

### Why tests use `cliutil.New`

The shared parser code exposes:

- `Options(cfg Config) []kong.Option`
- `New(app any, cfg Config, options ...kong.Option)`

This exists partly to make tests realistic without duplicating parser wiring. Tests can:

- reuse the same base parser configuration as production
- override writers
- override exit handling

## Runtime Behavior

At runtime, the current CLI behaves like this:

- binary name is `gsimp`
- description is `Starter CLI template`
- flags:
  - `--help`
  - `--config-file`
  - `--log-level`
  - `--version`
- commands:
  - `example`

`gsimp example` logs a single informational line including the resolved config file path.

## Exact Current Source Shape

The current template is intentionally small.

### `cmd/globals.go`

Contains only:

- package `cmd`
- `type Globals struct { ConfigFile string ... }`

### `cmd/gsimp/cli.go`

Contains only:

- package `main`
- imports Kong, shared `cmd`, `example`, and zerolog
- root `CLI` struct
- `GetLogLevel`
- `RunArgs`

### `cmd/gsimp/main.go`

Contains only:

- package `main`
- linker-overridable variables
- `main()`
- shared parser/runtime calls
- zerolog global formatting setup
- error logging and `os.Exit(1)`

### `cmd/gsimp/example/cmd.go`

Contains only:

- package `example`
- `type Command struct{}`
- `Run(log zerolog.Logger, g *cmd.Globals) error`

### `internal/cli/buildinfo.go`

Contains only:

- `type BuildInfo`

### `internal/cli/runner.go`

Contains:

- `Runner` interface
- `Config` struct
- `Options`
- `New`
- `Parse`
- `NewLogger`
- `Run`

## Design Decisions and Rationale

### Thin mains

Each binary should be mostly:

- build metadata
- description
- parser call
- logger setup
- run
- exit on error

This keeps new binaries cheap to add.

### Shared runtime package

`internal/cli` exists so future binaries can reuse:

- parser setup
- Kong option construction
- build metadata interpolation
- logger construction
- `ctx.Run(...)` glue

### Shared globals live in `cmd`, not `internal/cli`

This was a deliberate final choice.

Reason:

- `Globals` is not parser/runtime infrastructure
- it is app-level injected state
- subcommands need a shared importable type
- `cmd/globals.go` is simple and easy to customize when cloning the template

### Logger settings split

The logger setup is intentionally split:

- `internal/cli.NewLogger` configures the returned logger instance
- `cmd/gsimp/main.go` sets process-level zerolog global knobs:
  - `DurationFieldUnit`
  - `TimeFieldFormat`

That keeps process-level side effects explicit in `main`.

### Generic placeholder wording

The template is intentionally generic:

- app description: `Starter CLI template`
- shared config flag: `--config-file`
- placeholder command: `example`

This avoids dragging domain-specific wording into new projects created from the template.

## How to Recreate This Repo From Scratch

If starting from an empty directory, recreate the project in this order.

1. Create `go.mod`

- module: `github.com/blumsicle/gsimp`
- Go: `1.26.1`
- add direct dependencies:
  - Kong
  - zerolog
  - testify

2. Create directories

- `cmd/`
- `cmd/gsimp/`
- `cmd/gsimp/example/`
- `internal/cli/`

3. Create `cmd/globals.go`

- package `cmd`
- define `Globals` with the `ConfigFile` flag exactly as described above

4. Create `internal/cli/buildinfo.go`

- package `cli`
- define `BuildInfo`

5. Create `internal/cli/runner.go`

- define `Runner`, `Config`, `Options`, `New`, `Parse`, `NewLogger`, and `Run`
- match the current behavior described earlier

6. Create `cmd/gsimp/cli.go`

- package `main`
- embed `cmd.Globals`
- add `LogLevel`, `Version`, and `Example`
- implement `GetLogLevel()` and `RunArgs()`

7. Create `cmd/gsimp/example/cmd.go`

- package `example`
- define `Command`
- `Run(log zerolog.Logger, g *cmd.Globals) error`
- log `example command` with `config_file`

8. Create `cmd/gsimp/main.go`

- package `main`
- define `name`, `version`, `commit`
- set description to `Starter CLI template`
- call shared parser/runtime
- set zerolog global formatting in `main`
- log errors and exit

9. Create `Makefile`

- match the current target and variable behavior exactly
- especially:
  - app discovery from `cmd/*/main.go`
  - version/tag logic
  - `-X main.name/version/commit`
  - `rebuild`
  - `fmt` using `gofumpt`
  - `check`
  - `clean`

10. Create `README.md`

- keep it short and user-facing
- include rename steps, common commands, and layout

11. Create `cmd/gsimp/main_test.go`

- use `testify`
- add the three smoke tests described above

12. Run:

- `gofumpt -w .`
- `go mod tidy`
- `go test ./...`

## How to Adapt It for a New Tool

When cloning this template for a new project:

- change the module path
- rename `cmd/gsimp` to the new binary name
- update the `name` variable in the new `main.go`
- update the description string
- replace `cmd/globals.go` with flags appropriate for the new tool
- rename or delete the placeholder `example` command
- update or replace tests to reflect the new command tree

## What Not to Change Accidentally

These details are easy to break:

- Kong handler injection requires the exact same concrete `*cmd.Globals` type in both the root CLI and command handlers
- linker flags should target `main.name`, `main.version`, and `main.commit`
- `build` is incremental; `rebuild` is forceful
- `Parse` currently exits fatally on parse errors in production
- `help` and `version` tests expect an error because Kong’s exit hook is overridden during tests

## Current Expected Commands and Outputs

Expected commands:

- `make build`
- `make rebuild`
- `make install`
- `make test`
- `make vet`
- `make check`
- `make clean`

Expected CLI examples:

- `gsimp --help`
- `gsimp --version`
- `gsimp example`
- `gsimp --config-file ~/.config/starter/config.yaml example`

## Final Summary

This repository is a small, opinionated Go CLI template with:

- one example binary
- one example subcommand
- shared runtime wiring in `internal/cli`
- shared injected args in `cmd/globals.go`
- linker-injected build metadata
- a practical Makefile
- smoke tests for help/version/injection behavior

If another session needed to rebuild this project from scratch, reproducing the file structure and behaviors in this document should produce something extremely close to the current repo.
