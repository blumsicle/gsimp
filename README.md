# gsimp

Starter repo for Go CLI tools built with Kong and zerolog.

## Rename for a New Tool

1. Copy this repository to a new directory.
2. Update the module path in `go.mod`.
3. Rename `cmd/gsimp` to your binary name.
4. Update `name`, `version`, `commit`, and `Description` in `cmd/<your-binary>/main.go`.
5. Update shared flags in `cmd/globals.go` to match the new tool.
6. Rename or replace the placeholder `example` subcommand under `cmd/<your-binary>/example`.
7. Reinitialize Git history for the new repository.

## Common Commands

- `make build` builds versioned binaries into `bin/`.
- `make rebuild` forces a rebuild of versioned binaries.
- `make install` installs the current CLI with embedded build metadata.
- `make test` runs Go tests.
- `make vet` runs `go vet`.
- `make check` runs tests and vetting together.
- `make clean` removes built artifacts from `bin/`.

## Layout

- `cmd/<binary>` contains binary entrypoints and command wiring.
- `cmd/globals.go` contains shared injected arguments used by command handlers.
- `internal/cli` contains shared CLI runtime code for parsing, logging, and build metadata wiring.
- `Makefile` handles local build, install, and verification workflows.
