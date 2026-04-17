# bcli

Generate starter repos for Go CLI tools built with Kong and zerolog.

This repository is licensed under the MIT License. See [`LICENSE`](./LICENSE).

## Usage

`bcli create mycommand "CLI tool that does some cool stuff"` creates a
new starter project in `./mycommand`.

After generating the files, `bcli` also runs post steps to update
dependencies, tidy the module, initialize Git, and create an initial
commit. Each of those four post steps can be disabled through config or
`bcli create` flags.

Use `--inplace` to write the scaffold into the current directory instead
of creating a child directory. In-place generation is intended for empty
project directories and ignores configured `root_path` and
`project_dir_prefix`:

`bcli create --inplace mycommand "CLI tool that does some cool stuff"`

By default the generated module path uses just the project name. Set
`--git-location` or `git_location` in config if you want a fully
qualified module path such as `github.com/your-org/<name>`.

Use `--root-path` or `-r` to choose a different parent directory:

`bcli create --root-path ~/src mycommand "CLI tool that does some cool
stuff"`

Use `--project-dir-prefix`, `-p`, or `project_dir_prefix` in config to
prepend a string to the generated directory name without changing the
project name used inside the scaffold:

`bcli create -p local- mycommand "CLI tool that does some cool stuff"`

Use `--git-location` or `-g` to change the repository prefix used in the
generated `go.mod`:

`bcli create --git-location github.com/your-org mycommand "CLI tool
that does some cool stuff"`

Configuration can also be loaded from the file pointed to by
`--config-file`, which defaults to `~/.config/bcli/bcli.yaml`. The
generator currently supports these YAML keys:

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

Run `bcli config` to generate a config file with the current defaults,
then edit it to fit your environment.

Precedence is listed from lowest to highest:

1. built-in defaults
2. YAML config file
3. explicit CLI flags

Use `bcli config` to inspect the fully resolved config after defaults
and YAML file loading have been applied:

`bcli config`

`bcli config` preserves environment-variable references from the config
file. `bcli create` normalizes config before generation and currently
expands environment variables, `~`, and `~user` in `root_path`.

By default, `bcli config` writes YAML to stdout. Use `--output` or `-o`
to write it to a file instead:

`bcli config --output /tmp/bcli-resolved.yaml`

If the parent directories in the `--output` path do not exist,
`bcli config` creates them before writing the file.

Use `bcli completion <shell>` to print a shell completion script for
`zsh`, `bash`, or `fish`:

`bcli completion zsh`

The zsh output defines `_bcli`. To install it persistently in zsh, write
it into a directory on `fpath`, for example:

`mkdir -p ~/.zsh/completions && bcli completion zsh > ~/.zsh/completions/_bcli`

Use `--json` with `bcli create` to write structured creation metadata to
stdout. This is mainly intended for automation:

`bcli create --json mycommand "CLI tool that does some cool stuff"`

`bcli-mcp` runs a stdio MCP server that lets Codex create projects
through an installed `bcli` command. Its config file defaults to
`~/.config/bcli/bcli-mcp.yaml`. Install the binaries, then register the
server with Codex:

```sh
task install
codex mcp add bcli-project-generator -- bcli-mcp
```

For Codex sessions that should resume inside the new project, start
Codex in an empty target directory and let the MCP tool call
`bcli create --inplace`.

The generated project includes:

- a thin `main`
- shared app config in `internal/appconfig`
- shared CLI runtime code under `internal/cli`
- shared injected args in `cmd/globals.go` with a default config path
  at `~/.config/<project>/<project>.yaml`
- a `go.mod` whose `go` version matches the locally available Go
  toolchain used to run `bcli`
- a shell completion subcommand for `zsh`, `bash`, and `fish`
- an example subcommand
- a Taskfile with build, rebuild, test, vet, check, and clean tasks
- smoke tests for help, version, and command wiring

## Common Commands

- `bcli config` writes the resolved config as YAML to stdout or a file.
- `bcli completion <shell>` prints a completion script for `zsh`,
  `bash`, or `fish`.
- `bcli-mcp` starts the MCP server over stdio.
- `task build` builds versioned binaries into `bin/`.
- `task rebuild` removes built artifacts and rebuilds versioned binaries.
- `task install` installs the current CLI.
- `task lint` runs `golangci-lint`.
- `task test` runs Go tests.
- `task coverage` runs Go tests with coverage output.
- `task coverage-html` generates an HTML coverage report.
- `task vet` runs `go vet`.
- `task check` runs linting, tests, and vetting together.
- `task clean` removes built artifacts from `bin/`.

## Tooling

The `task fmt` task expects these tools to be installed locally:

- `gofumpt`
- `gci`
- `goimports`
- `golines`

Install Task with your package manager, for example `brew install
go-task/tap/go-task`.

## Layout

- `cmd/bcli` contains the generator binary entrypoint and commands.
- `cmd/bcli-mcp` contains the MCP server binary entrypoint.
- `internal/projectgen` contains the project generator.
- `internal/mcpserver` contains the MCP tool server and `bcli` shell-out
  adapter.
- `internal/bcliconfig`, `internal/poststep`, `cmd/globals.go`, and
  `internal/cli` define the generator/runtime pieces in this repo.
- `Taskfile.yml` handles local build, install, and verification workflows.
