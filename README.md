# bcli

Generate starter repos for Go CLI tools built with Kong and zerolog.

This repository is licensed under the MIT License. See [`LICENSE`](./LICENSE).

## Usage

`bcli create mycommand "CLI tool that does some cool stuff"` creates a new starter project in `./mycommand`.

After generating the files, `bcli` also runs post steps to update dependencies, tidy the module, initialize Git, and create an initial commit. Each of those four post steps can be disabled through config or `bcli create` flags.

By default the generated module path uses just the project name. Set `--git-location` or `git_location` in config if you want a fully qualified module path such as `github.com/your-org/<name>`.

Use `--root-path` or `-r` to choose a different parent directory:

`bcli create --root-path ~/src mycommand "CLI tool that does some cool stuff"`

Use `--project-dir-prefix`, `-p`, or `project_dir_prefix` in config to prepend a string to the generated directory name without changing the project name used inside the scaffold:

`bcli create -p local- mycommand "CLI tool that does some cool stuff"`

Use `--git-location` or `-g` to change the repository prefix used in the generated `go.mod`:

`bcli create --git-location github.com/your-org mycommand "CLI tool that does some cool stuff"`

Configuration can also be loaded from the file pointed to by `--config-file`. The generator currently supports these YAML keys:

```yaml
log_level: info
root_path: ~/src
project_dir_prefix: ""
git_location: ""
post_steps:
  go_get_update: true
  go_mod_tidy: true
  git_init: true
  git_commit: true
```

Run `bcli config` to generate a config file with the current defaults, then edit it to fit your environment.

Precedence is listed from lowest to highest:

1. built-in defaults
2. YAML config file
3. explicit CLI flags

Use `bcli config` to inspect the fully resolved config after defaults and YAML file loading have been applied:

`bcli config`

`bcli config` preserves environment-variable references from the config file. `bcli create` normalizes config before generation and currently expands environment variables in `root_path`.

By default, `bcli config` writes YAML to stdout. Use `--output` or `-o` to write it to a file instead:

`bcli config --output /tmp/bcli-resolved.yaml`

If the parent directories in the `--output` path do not exist, `bcli config` creates them before writing the file.

Use `bcli completion <shell>` to print a shell completion script for `zsh`, `bash`, or `fish`:

`bcli completion zsh`

The zsh output defines `_bcli`. To install it persistently in zsh, write it into a directory on `fpath`, for example:

`mkdir -p ~/.zsh/completions && bcli completion zsh > ~/.zsh/completions/_bcli`

The generated project includes:

- a thin `main`
- shared app config in `internal/appconfig`
- shared CLI runtime code under `internal/cli`
- shared injected args in `cmd/globals.go` with a default config path under `~/.config/<project>/`
- a `go.mod` whose `go` version matches the locally available Go toolchain used to run `bcli`
- a shell completion subcommand for `zsh`, `bash`, and `fish`
- an example subcommand
- a Makefile with build, rebuild, test, vet, check, and clean targets
- smoke tests for help, version, and command wiring

## Common Commands

- `bcli config` writes the resolved config as YAML to stdout or a file.
- `bcli completion <shell>` prints a completion script for `zsh`, `bash`, or `fish`.
- `make build` builds versioned binaries into `bin/`.
- `make rebuild` forces a rebuild of versioned binaries.
- `make install` installs the current CLI.
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

- `cmd/bcli` contains the generator binary entrypoint and commands.
- `internal/projectgen` contains the project generator.
- `internal/appconfig`, `internal/poststep`, `cmd/globals.go`, `internal/cli`, and `cmd/<binary>/example` define the generator/runtime pieces in this repo.
- `Makefile` handles local build, install, and verification workflows.
