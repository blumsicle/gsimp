package projectgen

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type exactDriftTemplate struct {
	name       string
	sourcePath string
	template   string
}

type divergentTemplate struct {
	name       string
	sourcePath string
	template   string
	reason     string
}

var exactDriftTemplates = []exactDriftTemplate{
	{
		name:       "gitignore",
		sourcePath: ".gitignore",
		template:   "templates/.gitignore.tmpl",
	},
	{
		name:       "makefile",
		sourcePath: "Makefile",
		template:   "templates/Makefile.tmpl",
	},
	{
		name:       "globals",
		sourcePath: "cmd/globals.go",
		template:   "templates/cmd/globals.go.tmpl",
	},
	{
		name:       "completion command",
		sourcePath: "cmd/bcli/completion/cmd.go",
		template:   "templates/cmd/__NAME__/completion/cmd.go.tmpl",
	},
	{
		name:       "completion command tests",
		sourcePath: "cmd/bcli/completion/cmd_test.go",
		template:   "templates/cmd/__NAME__/completion/cmd_test.go.tmpl",
	},
	{
		name:       "cli build info",
		sourcePath: "internal/cli/buildinfo.go",
		template:   "templates/internal/cli/buildinfo.go.tmpl",
	},
	{
		name:       "cli build info tests",
		sourcePath: "internal/cli/buildinfo_test.go",
		template:   "templates/internal/cli/buildinfo_test.go.tmpl",
	},
	{
		name:       "cli config loading",
		sourcePath: "internal/cli/config.go",
		template:   "templates/internal/cli/config.go.tmpl",
	},
	{
		name:       "cli config loading tests",
		sourcePath: "internal/cli/config_test.go",
		template:   "templates/internal/cli/config_test.go.tmpl",
	},
	{
		name:       "cli runner",
		sourcePath: "internal/cli/runner.go",
		template:   "templates/internal/cli/runner.go.tmpl",
	},
	{
		name:       "cli runner tests",
		sourcePath: "internal/cli/runner_test.go",
		template:   "templates/internal/cli/runner_test.go.tmpl",
	},
}

var intentionallyDivergentTemplates = []divergentTemplate{
	{
		name:       "readme",
		sourcePath: "README.md",
		template:   "templates/README.md.tmpl",
		reason:     "the repo README documents bcli as a generator; the scaffold README documents a generated CLI app",
	},
	{
		name:       "root cli",
		sourcePath: "cmd/bcli/cli.go",
		template:   "templates/cmd/__NAME__/cli.go.tmpl",
		reason:     "bcli exposes create while generated projects expose example",
	},
	{
		name:       "main command",
		sourcePath: "cmd/bcli/main.go",
		template:   "templates/cmd/__NAME__/main.go.tmpl",
		reason:     "bcli uses a command-specific config package while generated projects still use appconfig",
	},
	{
		name:       "config command",
		sourcePath: "cmd/bcli/config/cmd.go",
		template:   "templates/cmd/__NAME__/config/cmd.go.tmpl",
		reason:     "bcli uses a command-specific config package while generated projects still use appconfig",
	},
	{
		name:       "config command tests",
		sourcePath: "cmd/bcli/config/cmd_test.go",
		template:   "templates/cmd/__NAME__/config/cmd_test.go.tmpl",
		reason:     "bcli config command tests cover stdout and post-step fields while generated projects use a smaller scaffold config",
	},
	{
		name:       "cli test harness",
		sourcePath: "cmd/bcli/harness_test.go",
		template:   "templates/cmd/__NAME__/harness_test.go.tmpl",
		reason:     "bcli uses a command-specific config package while generated projects still use appconfig",
	},
	{
		name:       "go module",
		sourcePath: "go.mod",
		template:   "templates/go.mod.tmpl",
		reason:     "bcli includes MCP server dependencies while generated projects do not",
	},
	{
		name:       "root cli tests",
		sourcePath: "cmd/bcli/main_test.go",
		template:   "templates/cmd/__NAME__/main_test.go.tmpl",
		reason:     "bcli help includes create while generated project help includes example",
	},
	{
		name:       "config integration tests",
		sourcePath: "cmd/bcli/config_integration_test.go",
		template:   "templates/cmd/__NAME__/config_integration_test.go.tmpl",
		reason:     "bcli config integration tests include create flag precedence and post-step config assertions",
	},
	{
		name:       "example integration tests",
		sourcePath: "cmd/bcli/create/cmd_test.go",
		template:   "templates/cmd/__NAME__/example/cmd_test.go.tmpl",
		reason:     "bcli create tests exercise project generation while generated example tests exercise scaffold command wiring",
	},
	{
		name:       "bcli config schema",
		sourcePath: "internal/bcliconfig/config.go",
		template:   "templates/internal/appconfig/config.go.tmpl",
		reason:     "bcli uses a command-specific config package and includes post-step config while generated projects use a smaller scaffold config",
	},
	{
		name:       "command overrides",
		sourcePath: "internal/bcliconfig/create.go",
		template:   "templates/internal/appconfig/example.go.tmpl",
		reason:     "bcli create overrides and generated example overrides share a pattern but represent different command domains",
	},
	{
		name:       "example command",
		sourcePath: "cmd/bcli/create/cmd.go",
		template:   "templates/cmd/__NAME__/example/cmd.go.tmpl",
		reason:     "bcli create generates projects while generated example command demonstrates config injection",
	},
	{
		name:       "bcli config defaults tests",
		sourcePath: "internal/bcliconfig/config_test.go",
		template:   "templates/internal/appconfig/config_test.go.tmpl",
		reason:     "bcli uses a command-specific config package and includes post-step defaults while generated projects use a smaller scaffold config",
	},
	{
		name:       "bcli config create override tests",
		sourcePath: "internal/bcliconfig/create_test.go",
		template:   "templates/internal/appconfig/example_test.go.tmpl",
		reason:     "bcli create override tests and generated example override tests share a pattern but represent different command domains",
	},
	{
		name:       "bcli config load tests",
		sourcePath: "internal/bcliconfig/load_test.go",
		template:   "templates/internal/appconfig/load_test.go.tmpl",
		reason:     "bcli uses a command-specific config package and includes post-step config while generated projects use a smaller scaffold config",
	},
	{
		name:       "bcli config normalize tests",
		sourcePath: "internal/bcliconfig/normalize_test.go",
		template:   "templates/internal/appconfig/normalize_test.go.tmpl",
		reason:     "bcli uses a command-specific config package while generated normalize tests use project-specific environment variable names",
	},
	{
		name:       "bcli config load",
		sourcePath: "internal/bcliconfig/load.go",
		template:   "templates/internal/appconfig/load.go.tmpl",
		reason:     "bcli uses a command-specific config package while generated projects still use appconfig",
	},
	{
		name:       "bcli config normalize",
		sourcePath: "internal/bcliconfig/normalize.go",
		template:   "templates/internal/appconfig/normalize.go.tmpl",
		reason:     "bcli uses a command-specific config package while generated projects still use appconfig",
	},
	{
		name:       "bcli config root overrides",
		sourcePath: "internal/bcliconfig/root.go",
		template:   "templates/internal/appconfig/root.go.tmpl",
		reason:     "bcli uses a command-specific config package while generated projects still use appconfig",
	},
	{
		name:       "bcli config root override tests",
		sourcePath: "internal/bcliconfig/root_test.go",
		template:   "templates/internal/appconfig/root_test.go.tmpl",
		reason:     "bcli uses a command-specific config package while generated projects still use appconfig",
	},
}

func TestTemplatesMatchCanonicalSourcesWithBCLIData(t *testing.T) {
	for _, tt := range exactDriftTemplates {
		t.Run(tt.name, func(t *testing.T) {
			got := renderTemplateWithBCLIData(t, tt.template)

			want, err := os.ReadFile(repoPath(tt.sourcePath))
			require.NoError(t, err)

			assert.Equal(t, string(want), got)
		})
	}
}

func TestIntentionallyDivergentTemplatesAreDocumented(t *testing.T) {
	for _, tt := range intentionallyDivergentTemplates {
		t.Run(tt.name, func(t *testing.T) {
			assert.FileExists(t, repoPath(tt.sourcePath))
			assertTemplateExists(t, tt.template)
			assert.NotEmpty(t, tt.reason)
		})
	}
}

func TestAllTemplatesHaveDriftClassification(t *testing.T) {
	classified := map[string]bool{}
	for _, tt := range exactDriftTemplates {
		classified[tt.template] = true
	}
	for _, tt := range intentionallyDivergentTemplates {
		classified[tt.template] = true
	}

	var unclassified []string
	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		if !classified[path] {
			unclassified = append(unclassified, path)
		}
		return nil
	})
	require.NoError(t, err)

	sort.Strings(unclassified)
	assert.Empty(t, unclassified)
}

func renderTemplateWithBCLIData(t *testing.T, templatePath string) string {
	t.Helper()

	got, err := New(zerolog.Nop()).renderTemplate(
		templatePath,
		templateData{
			Name:        "bcli",
			Description: "Generate starter Go CLI projects",
			ModulePath:  "github.com/blumsicle/bcli",
			GoVersion:   repoGoVersion(t),
		},
	)
	require.NoError(t, err)

	return got
}

func assertTemplateExists(t *testing.T, templatePath string) {
	t.Helper()

	_, err := fs.Stat(templateFS, templatePath)
	require.NoError(t, err)
}

func repoPath(relativePath string) string {
	return filepath.Join("..", "..", filepath.FromSlash(relativePath))
}

func repoGoVersion(t *testing.T) string {
	t.Helper()

	goMod, err := os.ReadFile(repoPath("go.mod"))
	require.NoError(t, err)

	for _, line := range strings.Split(string(goMod), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "go" {
			return fields[1]
		}
	}

	require.Fail(t, "go.mod does not contain a go directive")
	return ""
}
