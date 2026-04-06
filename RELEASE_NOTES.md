# Release Notes

## v0.1.0 - 2026-04-06

First public release of `gsimp`.

`gsimp` is now a generator CLI for bootstrapping Go command-line applications built with Kong and zerolog. This initial release focuses on a clean generated project layout, explicit configuration, and a practical post-generation workflow.

### Highlights

- Added `gsimp create <name> <description>` to generate a new Go CLI starter project.
- Generated projects include a thin `main`, root command wiring, shared globals, an example subcommand, app config loading, runtime helpers, tests, a `Makefile`, and an example YAML config file.
- Added typed app configuration with defaults, YAML loading, environment expansion, and clear precedence between defaults, config files, and CLI flags.
- Added post-generation steps for dependency update, module tidy, Git initialization, and initial commit creation.
- Split Git initialization and initial commit into separate post steps and added dependency-aware planning so `git commit` is skipped automatically when `git init` is disabled.
- Added per-step configuration and CLI flags to disable any of the four post steps.
- Improved structured logging across the generator with subsystem-specific logger names for command handling, project generation, and post-step execution.

### Generated Project Behavior

Generated projects currently include:

- `cmd/<project>` for the binary entrypoint and command tree
- `cmd/globals.go` for shared CLI globals
- `internal/appconfig` for typed config defaults, YAML loading, and tests
- `internal/cli` for parser, build info, and logger helpers
- `cmd/<project>/example` as a starter subcommand
- a version-aware `Makefile`
- a `<project>.yaml` example config file

By default, generation also runs:

- `go get -u ./...`
- `go mod tidy`
- `git init`
- `git add .`
- `git commit -m "Initial commit"`

### Notes

- This is an early `0.x` release. The generator is usable, but generated output and internal structure may still evolve before a `v1.0.0` stability commitment.
