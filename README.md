# gsimp

Generate starter repos for Go CLI tools built with Kong and zerolog.

## Usage

`gsimp create mycommand "CLI tool that does some cool stuff"` creates a new starter project in `./mycommand`.

By default the generated module path uses just the project name. Set `--git-location` or `git_location` in config if you want a fully qualified module path such as `github.com/your-org/<name>`.

Use `--root-path` or `-r` to choose a different parent directory:

`gsimp create --root-path ~/src mycommand "CLI tool that does some cool stuff"`

Use `--git-location` or `-g` to change the repository prefix used in the generated `go.mod`:

`gsimp create --git-location github.com/your-org mycommand "CLI tool that does some cool stuff"`

Configuration can also be loaded from the file pointed to by `--config-file`. The generator currently supports these YAML keys:

```yaml
log_level: info
root_path: ~/src
git_location: ""
```

See [`gsimp.yaml`](./gsimp.yaml) for a concrete example config file with comments.

Precedence is:

1. built-in defaults
2. YAML config file
3. explicit CLI flags

The generated project includes:

- a thin `main`
- shared app config in `cmd/config.go` with defaults, YAML loading, and env var expansion
- shared CLI runtime code under `internal/cli`
- shared injected args in `cmd/globals.go` with a default config path under `~/.config/<project>/`
- a `<project>.yaml` example config file
- an example subcommand
- a Makefile with build, rebuild, test, vet, check, and clean targets
- smoke tests for help, version, and command wiring

## Common Commands

- `make build` builds versioned binaries into `bin/`.
- `make rebuild` forces a rebuild of versioned binaries.
- `make install` installs the current CLI with embedded build metadata.
- `make test` runs Go tests.
- `make coverage` runs Go tests with coverage output.
- `make coverage-html` generates an HTML coverage report.
- `make vet` runs `go vet`.
- `make check` runs tests and vetting together.
- `make clean` removes built artifacts from `bin/`.

## Tooling

The `make fmt` target expects these tools to be installed locally:

- `gofumpt`
- `goimports`
- `golines`

## Layout

- `cmd/gsimp` contains the generator binary entrypoint and commands.
- `internal/projectgen` contains the project generator.
- `cmd/config.go`, `cmd/globals.go`, `internal/cli`, and `cmd/<binary>/example` define the starter project that gets generated.
- `Makefile` handles local build, install, and verification workflows.
