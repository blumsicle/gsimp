# Feature Roadmap

Living document for product-facing `bcli` feature ideas.

Status values:

- `candidate`
- `planned`
- `in progress`
- `done`

## Backlog

### 1. Add `bcli doctor`

- Status: `candidate`
- Priority: high
- Summary:
  Add a diagnostic command that checks whether the local environment is
  ready to generate, build, test, and maintain `bcli` projects.
- Candidate checks:
  - `go` is installed and reports a version.
  - `git` is installed.
  - formatter tools used by generated Makefiles are installed:
    `gofumpt`, `gci`, `goimports`, and `golines`.
  - `golangci-lint` is installed.
  - configured `root_path` normalizes and is usable.
  - configured `git_location` looks like a valid module prefix when set.
- Notes:
  - This is the recommended first feature because it supports both
    repository development and generated project workflows without
    increasing template complexity.

### 2. Add `create --dry-run`

- Status: `candidate`
- Priority: high
- Summary:
  Let users resolve config and create inputs without writing files or
  running post-generation steps.
- Candidate output:
  - target path
  - module path
  - project name
  - description
  - files that would be written
  - post steps that would run
- Notes:
  - This can be a simpler first step toward a richer preview command.

### 3. Validate Project Names And Module Inputs

- Status: `candidate`
- Priority: high
- Summary:
  Add stricter validation before generation so invalid project names or
  module paths fail early with clear messages.
- Candidate checks:
  - project name is non-empty.
  - project name does not contain path separators.
  - project name avoids spaces and characters that produce invalid Go
    package/import paths.
  - `git_location` plus project name forms an expected module path.
  - target path does not already exist.
- Notes:
  - This should integrate with `create` and any future `validate` or
    `preview` command.

### 4. Add `--no-post-steps`

- Status: `candidate`
- Priority: medium
- Summary:
  Add a convenience flag to disable all post-generation steps at once.
- Notes:
  - This complements the existing individual flags:
    `--no-go-get-update`, `--no-go-mod-tidy`, `--no-git-init`, and
    `--no-git-commit`.
  - Useful for tests, quick scaffolding, and CI-style dry generation.

### 5. Add `bcli preview`

- Status: `candidate`
- Priority: medium
- Summary:
  Show generated output without writing a project.
- Candidate modes:
  - `bcli preview <name> <description>` prints a file tree.
  - `--tree` prints only generated paths.
  - `--file <path>` prints one rendered file.
  - `--diff` can later compare generated content against an existing
    target directory.
- Notes:
  - This may share implementation with `create --dry-run`, but the user
    intent is different: `dry-run` validates create behavior, while
    `preview` inspects rendered output.

### 6. Add Config Initialization Flow

- Status: `candidate`
- Priority: medium
- Summary:
  Make first-run config setup clearer.
- Candidate shapes:
  - `bcli config init`
  - `bcli config --output ~/.config/bcli/config.yaml`
  - `--force` to allow overwriting an existing config file.
- Notes:
  - Current `bcli config` already writes resolved YAML. This feature
    would make the default setup path more explicit.

### 7. Add `bcli templates`

- Status: `candidate`
- Priority: low
- Summary:
  Add maintainer-oriented template inspection commands.
- Candidate subcommands:
  - `bcli templates list`
  - `bcli templates show <path>`
  - `bcli templates drift`
- Notes:
  - Internal drift tests already protect templates. A CLI command would
    mainly help maintainers inspect embedded scaffolds outside tests.

### 8. Add Existing Directory Controls

- Status: `candidate`
- Priority: low
- Summary:
  Give users safer options when a target directory already exists.
- Candidate options:
  - `--allow-existing-empty-dir`
  - `--force` or `--overwrite` for explicit overwrites
- Notes:
  - Prefer `--allow-existing-empty-dir` before broader overwrite
    behavior. Overwriting scaffold files has a higher risk of deleting
    user work.

### 9. Add Generated Project Options

- Status: `candidate`
- Priority: low
- Summary:
  Let users choose optional scaffold pieces during generation.
- Candidate options:
  - `--no-example`
  - `--no-completion`
  - `--no-config-command`
  - `--license mit`
  - `--github-actions`
  - `--goreleaser`
- Notes:
  - Defer until the core workflow is stable. These options increase
    template branching and test matrix size.

## Suggested Order

1. Add `bcli doctor`.
2. Add `create --dry-run`.
3. Validate project names and module inputs.
4. Add `--no-post-steps`.
5. Add `bcli preview`.
6. Add config initialization flow.
7. Add `bcli templates`.
8. Add existing directory controls.
9. Add generated project options.
