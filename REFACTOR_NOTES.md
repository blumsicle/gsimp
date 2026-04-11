# Refactor Notes

Living document for maintainability cleanup work.

Status values:

- `pending`
- `in progress`
- `done`

## Backlog

### 1. Reduce repo/template duplication

- Status: `pending`
- Priority: high
- Summary:
  The real CLI/runtime code and the generated template code mirror each
  other closely, which makes most feature changes touch both the real
  files and the templates.
- Key references:
  - [internal/projectgen/generator.go](/Users/blumsicle/src/go/bcli/internal/projectgen/generator.go)
  - [internal/cli/runner.go](/Users/blumsicle/src/go/bcli/internal/cli/runner.go)
  - [internal/projectgen/templates/internal/cli/runner.go.tmpl](/Users/blumsicle/src/go/bcli/internal/projectgen/templates/internal/cli/runner.go.tmpl)
  - [cmd/bcli/main_test.go](/Users/blumsicle/src/go/bcli/cmd/bcli/main_test.go)
  - [internal/projectgen/templates/cmd/__NAME__/main_test.go.tmpl](/Users/blumsicle/src/go/bcli/internal/projectgen/templates/cmd/__NAME__/main_test.go.tmpl)
- Progress notes:
  - 2026-04-09: Identified as the highest-value refactor because it
    reduces change cost across repo code and generated scaffold code.
  - 2026-04-11: Tilde expansion for config paths required parallel
    changes in real appconfig code/tests and generated appconfig
    templates/tests, confirming appconfig as another duplication
    hotspot.

### 2. Split `Generator.Generate` into smaller private steps

- Status: `done`
- Priority: high
- Summary:
  `Generator.Generate` currently validates input, resolves paths, builds
  template data, renders files, writes files, and runs post-steps in one
  method.
- Key references:
  - [internal/projectgen/generator.go:64](/Users/blumsicle/src/go/bcli/internal/projectgen/generator.go#L64)
- Progress notes:
  - 2026-04-09: Candidate extraction points include config validation,
    target-path resolution, template rendering and writing, and post-step
    execution.
  - 2026-04-11: Split `Generate` into orchestration plus private helpers
    for config validation, generation planning, target directory setup,
    template rendering, and post-step execution.

### 3. Centralize config override application

- Status: `done`
- Priority: medium
- Summary:
  Root-level and command-level config precedence is currently applied in
  CLI handlers instead of being owned in one config-focused place.
- Key references:
  - [cmd/bcli/cli.go:24](/Users/blumsicle/src/go/bcli/cmd/bcli/cli.go#L24)
  - [cmd/bcli/create/cmd.go:28](/Users/blumsicle/src/go/bcli/cmd/bcli/create/cmd.go#L28)
  - [internal/appconfig/config.go](/Users/blumsicle/src/go/bcli/internal/appconfig/config.go)
- Progress notes:
  - 2026-04-09: Good candidate for explicit `Apply...Overrides` helpers
    on config or config-adjacent types.
  - 2026-04-11: Added appconfig-owned root and create override helpers
    so command handlers adapt flag values but no longer own config
    precedence mutation.

### 4. Simplify post-step definitions and planning

- Status: `done`
- Priority: medium
- Summary:
  The planner rebuilds the definition table on each call and uses small
  one-off closures for straightforward enabled-state checks.
- Key references:
  - [internal/poststep/planner.go:46](/Users/blumsicle/src/go/bcli/internal/poststep/planner.go#L46)
- Progress notes:
  - 2026-04-09: A more static definition table plus lightweight step
    factories would likely be easier to scan and maintain.
  - 2026-04-11: Moved planner definitions to a static private table and
    replaced per-call anonymous enabled closures with small spec methods
    for enabled-state checks and step construction.

### 5. Factor repetitive command-backed post-step implementations

- Status: `done`
- Priority: medium
- Summary:
  Several post-step implementations only differ by step name, log
  message, and command arguments.
- Key references:
  - [internal/poststep/go_get_update.go](/Users/blumsicle/src/go/bcli/internal/poststep/go_get_update.go)
  - [internal/poststep/go_mod_tidy.go](/Users/blumsicle/src/go/bcli/internal/poststep/go_mod_tidy.go)
  - [internal/poststep/git_init.go](/Users/blumsicle/src/go/bcli/internal/poststep/git_init.go)
  - [internal/poststep/git_commit.go](/Users/blumsicle/src/go/bcli/internal/poststep/git_commit.go)
- Progress notes:
  - 2026-04-09: `git commit` likely remains custom, but the
    single-command steps could be collapsed behind a small helper type.
  - 2026-04-11: Added shared private command post-step specs for the
    single-command steps while keeping `git commit` custom.

### 6. Introduce a small CLI test harness helper

- Status: `pending`
- Priority: low
- Summary:
  Parser, stdout, stderr, and exit-code setup is repeated across the
  repo CLI smoke tests and again in the generated template tests.
- Key references:
  - [cmd/bcli/main_test.go:29](/Users/blumsicle/src/go/bcli/cmd/bcli/main_test.go#L29)
  - [internal/projectgen/templates/cmd/__NAME__/main_test.go.tmpl](/Users/blumsicle/src/go/bcli/internal/projectgen/templates/cmd/__NAME__/main_test.go.tmpl)
- Progress notes:
  - 2026-04-09: This is lower priority than the structural
    production-code refactors, but it would reduce test noise.

## Suggested Order

1. Reduce repo/template duplication.
2. Split `Generator.Generate`.
3. Centralize config override application.
4. Simplify post-step definitions and implementations.
5. Add a CLI test harness helper.
