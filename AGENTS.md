# Repository Guidelines

## Project Structure & Module Organization

`gsimp` is a Go CLI generator. The main entrypoint lives in `cmd/gsimp`, with the `create` subcommand under `cmd/gsimp/create`. Shared CLI globals are in `cmd/globals.go`.

Core packages live under `internal/`:

- `internal/appconfig`: config defaults, YAML loading, and config tests
- `internal/cli`: shared runtime helpers such as runner and build info
- `internal/projectgen`: template rendering and generation orchestration
- `internal/poststep`: side-effecting post-generation steps such as `go mod tidy` and `git init`

Embedded scaffold templates are in `internal/projectgen/templates`. When changing generated project behavior, update both generator code and the matching `.tmpl` files.

## Build, Test, and Development Commands

Use `make` targets as the primary interface:

- `make install`: install the local CLI from `cmd/*/main.go`
- `make build`: build versioned binaries into `bin/`
- `make fmt`: run `gofumpt`, `goimports`, and `golines`
- `make test`: run `go test ./...`
- `make vet`: run `go vet ./...`
- `golangci-lint run ./...`: run the full lint suite expected by the repo

Standard verification flow after changes: `golangci-lint run ./...`, `make fmt`, then `make test`.

Run `golangci-lint run ./...` outside the sandbox for this repository. Linting relies on normal Go build cache access to analyze the module correctly.
Run `make test` outside the sandbox for this repository. The test suite invokes `go` and `git` in temporary directories and relies on normal Go build cache access.

## Coding Style & Naming Conventions

Follow standard Go formatting and keep code `gofumpt`-clean. Use tabs as emitted by Go tooling. Prefer package names that are short and descriptive (`poststep`, `projectgen`), exported identifiers in `CamelCase`, and tests in `*_test.go`.

For any new code, add Go doc comments to every exported type, function, method, constant, and variable. Keep comments concise and aligned with GoDoc conventions.

Template paths use the `__NAME__` placeholder for generated project names. Keep template file names aligned with their generated output paths.

## Testing Guidelines

Tests use Go’s built-in `testing` package. Keep tests adjacent to the code they cover, for example `internal/projectgen/generator_test.go`. Name tests with `Test...` and favor table-driven cases for config and generation behavior.

Some post-step tests invoke `go` and `git` in temporary directories. Preserve that behavior when refactoring.

## Commit & Pull Request Guidelines

Recent commits use short, imperative subjects such as `Add project generation post steps` and `Refactor generator config and post step packages`. Keep commit titles concise and descriptive.

Pull requests should explain the user-visible change, note affected packages or templates, and include the verification commands you ran. If template output changes, mention the generated scaffold impact explicitly.
