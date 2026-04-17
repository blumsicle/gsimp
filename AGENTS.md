# Repository Guidelines

## Project Structure & Module Organization

`bcli` is a Go CLI generator. The main entrypoint lives in `cmd/bcli`,
with the `create` subcommand under `cmd/bcli/create`. Shared CLI globals
are in `cmd/globals.go`.

Repository identity:

- GitHub repo: `github.com/blumsicle/bcli`
- Go module: `github.com/blumsicle/bcli`
- Current local checkout path: `/Users/blumsicle/src/go/bcli`

Core packages live under `internal/`:

- `internal/appconfig`: config defaults, YAML loading, and config tests
- `internal/cli`: shared runtime helpers such as runner and build info
- `internal/projectgen`: template rendering and generation orchestration
- `internal/poststep`: side-effecting post-generation steps such as
  `go mod tidy` and `git init`

Embedded scaffold templates are in `internal/projectgen/templates`. When
changing generated project behavior, update both generator code and the
matching `.tmpl` files.

## Build, Test, and Development Commands

Use `task` tasks as the primary interface:

- `task install`: install the local CLI from `cmd/*/main.go`
- `task build`: build versioned binaries into `bin/`
- `task fmt`: run `gofumpt`, `gci`, `goimports`, and `golines`
- `task lint`: run `golangci-lint run ./...`
- `task test`: run `go test ./...`
- `task vet`: run `go vet ./...`
- `task check`: run `task lint`, `task vet`, and `task test`
- `golangci-lint run ./...`: run the full lint suite expected by the repo

Standard verification flow after changes: `task fmt`, then `task check`.

Run `task check` outside the sandbox for this repository. The
verification flow invokes `go` and `git` in temporary directories and
relies on normal Go build cache access.
After code changes, also review user-facing documentation files such as
`README.md`, sample config files, release notes, and generated template
docs for any needed updates, and make those documentation changes in the
same pass when applicable.

## Coding Style & Naming Conventions

Follow standard Go formatting and keep code `gofumpt`-clean. Use tabs as
emitted by Go tooling. Prefer package names that are short and
descriptive (`poststep`, `projectgen`), exported identifiers in
`CamelCase`, and tests in `*_test.go`.

For any new code, add Go doc comments to every exported type, function,
method, constant, and variable. Keep comments concise and aligned with
GoDoc conventions.

Template paths use the `__NAME__` placeholder for generated project
names. Keep template file names aligned with their generated output
paths.

## Testing Guidelines

Tests use Go’s built-in `testing` package. Keep tests adjacent to the
code they cover, for example `internal/projectgen/generator_test.go`.
Name tests with `Test...` and favor table-driven cases for config and
generation behavior.

Some post-step tests invoke `go` and `git` in temporary directories.
Preserve that behavior when refactoring.

## Release Process

When asked to create a release, write or update a descriptive
`RELEASE_NOTES.md` manually rather than auto-generating a minimal
changelog. Then run `task fmt` and `task check` before tagging. Use
`scripts/release.sh --no-ask vX.Y.Z` after `RELEASE_NOTES.md` is
already updated in the worktree; the script stages `RELEASE_NOTES.md`,
creates the release commit, creates an annotated tag, and pushes both
the branch and tag to `origin`.

When run manually without `--no-ask`, `scripts/release.sh` prompts for
confirmation that `RELEASE_NOTES.md` was updated before proceeding.

The GitHub Actions workflow at `.github/workflows/release.yml`
publishes releases for pushed `v*` tags via GoReleaser using
`.goreleaser.yaml`. It runs `task fmt` and `task check`, builds the
`darwin/arm64` `bcli` binary, and publishes the GitHub release using
`RELEASE_NOTES.md` from the tagged commit as the release body.
GitHub repository settings must allow workflow `contents: write`
permissions for the release job to create GitHub releases successfully,
and the workflow must provide `GITHUB_TOKEN` to GoReleaser.

## Commit & Pull Request Guidelines

Recent commits use short, imperative subjects such as `Add project
generation post steps` and `Refactor generator config and post step
packages`. Keep commit titles concise and descriptive.

Pull requests should explain the user-visible change, note affected
packages or templates, and include the verification commands you ran. If
template output changes, mention the generated scaffold impact
explicitly.
