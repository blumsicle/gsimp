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
		name:       "go module",
		sourcePath: "go.mod",
		template:   "templates/go.mod.tmpl",
	},
	{
		name:       "main command",
		sourcePath: "cmd/bcli/main.go",
		template:   "templates/cmd/__NAME__/main.go.tmpl",
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
		name:       "config command",
		sourcePath: "cmd/bcli/config/cmd.go",
		template:   "templates/cmd/__NAME__/config/cmd.go.tmpl",
	},
	{
		name:       "cli test harness",
		sourcePath: "cmd/bcli/harness_test.go",
		template:   "templates/cmd/__NAME__/harness_test.go.tmpl",
	},
	{
		name:       "completion command tests",
		sourcePath: "cmd/bcli/completion/cmd_test.go",
		template:   "templates/cmd/__NAME__/completion/cmd_test.go.tmpl",
	},
	{
		name:       "appconfig load",
		sourcePath: "internal/appconfig/load.go",
		template:   "templates/internal/appconfig/load.go.tmpl",
	},
	{
		name:       "appconfig normalize",
		sourcePath: "internal/appconfig/normalize.go",
		template:   "templates/internal/appconfig/normalize.go.tmpl",
	},
	{
		name:       "appconfig root overrides",
		sourcePath: "internal/appconfig/root.go",
		template:   "templates/internal/appconfig/root.go.tmpl",
	},
	{
		name:       "appconfig root override tests",
		sourcePath: "internal/appconfig/root_test.go",
		template:   "templates/internal/appconfig/root_test.go.tmpl",
	},
	{
		name:       "cli build info",
		sourcePath: "internal/cli/buildinfo.go",
		template:   "templates/internal/cli/buildinfo.go.tmpl",
	},
	{
		name:       "cli runner",
		sourcePath: "internal/cli/runner.go",
		template:   "templates/internal/cli/runner.go.tmpl",
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
		name:       "root cli tests",
		sourcePath: "cmd/bcli/main_test.go",
		template:   "templates/cmd/__NAME__/main_test.go.tmpl",
		reason:     "bcli help includes create while generated project help includes example",
	},
	{
		name:       "config integration tests",
		sourcePath: "cmd/bcli/config_integration_test.go",
		template:   "templates/cmd/__NAME__/config/cmd_test.go.tmpl",
		reason:     "bcli config integration tests include create flag precedence and post-step config assertions",
	},
	{
		name:       "example integration tests",
		sourcePath: "cmd/bcli/create/cmd_test.go",
		template:   "templates/cmd/__NAME__/example/cmd_test.go.tmpl",
		reason:     "bcli create tests exercise project generation while generated example tests exercise scaffold command wiring",
	},
	{
		name:       "appconfig schema",
		sourcePath: "internal/appconfig/config.go",
		template:   "templates/internal/appconfig/config.go.tmpl",
		reason:     "bcli includes post-step config while generated projects use a smaller scaffold config",
	},
	{
		name:       "command overrides",
		sourcePath: "internal/appconfig/create.go",
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
		name:       "appconfig defaults tests",
		sourcePath: "internal/appconfig/config_test.go",
		template:   "templates/internal/appconfig/config_test.go.tmpl",
		reason:     "bcli config tests include post-step defaults while generated projects use a smaller scaffold config",
	},
	{
		name:       "appconfig example override tests",
		sourcePath: "internal/appconfig/create_test.go",
		template:   "templates/internal/appconfig/example_test.go.tmpl",
		reason:     "bcli create override tests and generated example override tests share a pattern but represent different command domains",
	},
	{
		name:       "appconfig load tests",
		sourcePath: "internal/appconfig/load_test.go",
		template:   "templates/internal/appconfig/load_test.go.tmpl",
		reason:     "bcli load tests include post-step config while generated projects use a smaller scaffold config",
	},
	{
		name:       "appconfig normalize tests",
		sourcePath: "internal/appconfig/normalize_test.go",
		template:   "templates/internal/appconfig/normalize_test.go.tmpl",
		reason:     "generated normalize tests use project-specific environment variable names",
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
