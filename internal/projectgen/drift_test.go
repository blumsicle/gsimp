package projectgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplatesMatchCanonicalSourcesWithBCLIData(t *testing.T) {
	tests := []struct {
		name       string
		sourcePath string
		template   string
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderTemplateWithBCLIData(t, tt.template)

			want, err := os.ReadFile(repoPath(tt.sourcePath))
			require.NoError(t, err)

			assert.Equal(t, string(want), got)
		})
	}
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
