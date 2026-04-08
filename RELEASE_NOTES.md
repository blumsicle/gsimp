# Release Notes

## v0.2.3 - 2026-04-07

Adds first-party shell completion support and propagates the current CLI/documentation behavior into generated projects.

### Highlights

- Added an MIT `LICENSE` file to the main `bcli` repository.
- Added `bcli completion <shell>` support for `zsh`, `bash`, and `fish` using `github.com/miekg/king`.
- Zsh completion output now emits a real `_bcli` completion definition suitable for installation on `fpath`.
- Added path-aware completion tags for config file and output flags in the repository CLI.
- Propagated the `completion` subcommand, dependency updates, tests, and README guidance into generated projects.
- Updated repository guidance to use `make check` as the standard verification target.

### Notes

- The MIT license still applies to this repository only; generated projects do not receive a license file automatically.
- This release intentionally replaces the previously published `v0.2.3` tag and GitHub release contents.
- This remains a `0.x` release, so generated output and command surface may still evolve before a `v1.0.0` stability commitment.

## v0.2.2 - 2026-04-07

Adds more flexible project output path handling in `create` and improves file output behavior in `config`, along with a generated `go.mod` update that follows the local Go toolchain.

### Highlights

- Added `project_dir_prefix` config support for `bcli create`.
- Added `--project-dir-prefix` and short flag `-p` to prepend a string to the generated project directory name without changing the project name used inside the scaffold.
- Updated project generation so the directory prefix only affects the target path, not the generated module path, package names, or binary name.
- Updated `bcli config --output` to create missing parent directories before writing the output file.
- Updated sample config and repository documentation to describe the new create-path option and config output behavior.

### Notes

- The default `project_dir_prefix` is an empty string, so existing `bcli create` behavior is unchanged unless the option is set.
- This remains a `0.x` release, so generated output and command surface may still evolve before a `v1.0.0` stability commitment.

## v0.2.1 - 2026-04-07

Refines build metadata reporting so installed binaries and generated projects report more useful version information without relying on linker flags.

### Highlights

- Added runtime build-info fallback in `internal/cli` so `go install github.com/blumsicle/bcli/cmd/bcli@latest` can report module and VCS metadata when available.
- Shortened runtime fallback commit display to a 7-character revision.
- Simplified binary entrypoints and generated templates so build metadata resolution only requires an explicit CLI name.
- Removed Makefile linker flag injection for name, version, and commit in favor of runtime build metadata resolution.
- Updated generator documentation and repository guidance to match the current metadata and documentation-update workflow.

### Notes

- Local builds without embedded module or VCS metadata still fall back to `dev` and `unknown`.
- This remains a `0.x` release, so generated output and command surface may still evolve before a `v1.0.0` stability commitment.

## v0.2.0 - 2026-04-07

Adds a resolved-config inspection command and finishes the remaining rename cleanup after the `bcli` transition.

### Highlights

- Added `bcli config` to print the fully resolved config as YAML after defaults and config-file loading are applied.
- Added `--output` / `-o` to `bcli config` so resolved YAML can be written to a file instead of stdout.
- Added command-level debug and info logging for the `config` subcommand.
- Finished the remaining `gsimp` to `bcli` rename cleanup in config-related tests.
- Updated repository docs to describe the new command and current generated-template behavior.

### Notes

- The `config` command outputs resolved YAML from the in-memory config model; it does not preserve comments or blank lines from the source config file.
- This remains a `0.x` release, so generated output and command surface may still evolve before a `v1.0.0` stability commitment.

## v0.1.0 - 2026-04-06

First public release of `bcli`.

`bcli` is now a generator CLI for bootstrapping Go command-line applications built with Kong and zerolog. This initial release focuses on a clean generated project layout, explicit configuration, and a practical post-generation workflow.

### Highlights

- Added `bcli create <name> <description>` to generate a new Go CLI starter project.
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
