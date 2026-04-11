package projectgen

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/blumsicle/bcli/internal/poststep"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCreatesStarterProject(t *testing.T) {
	rootPath := t.TempDir()

	targetPath, err := New(zerolog.Nop()).Generate(context.Background(), Config{
		Name:             "mycommand",
		Description:      "CLI tool that does some cool stuff",
		GitLocation:      "github.com/blumsicle",
		ProjectDirPrefix: "generated-",
		RootPath:         rootPath,
	})
	require.NoError(t, err)

	require.Equal(t, filepath.Join(rootPath, "generated-mycommand"), targetPath)

	goMod, err := os.ReadFile(filepath.Join(targetPath, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(goMod), "module github.com/blumsicle/mycommand")
	assert.Contains(t, string(goMod), "go "+currentGoVersion())

	mainGo, err := os.ReadFile(filepath.Join(targetPath, "cmd", "mycommand", "main.go"))
	require.NoError(t, err)
	assert.Contains(t, string(mainGo), "Description: \"CLI tool that does some cool stuff\"")
	assert.Contains(t, string(mainGo), "const name = \"mycommand\"")
	assert.Contains(
		t,
		string(mainGo),
		"BuildInfo:   cliutil.ResolveBuildInfo(name),",
	)

	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "config.go"))
	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "example.go"))
	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "load.go"))
	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "normalize.go"))
	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "root.go"))
	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "config_test.go"))

	loadTest, err := os.ReadFile(
		filepath.Join(targetPath, "internal", "appconfig", "load_test.go"),
	)
	require.NoError(t, err)
	assert.Contains(
		t,
		string(loadTest),
		`[]byte("root_path: $MYCOMMAND_TEST_HOME/src\ngit_location: github.com/acme\n")`,
	)

	normalizeTest, err := os.ReadFile(
		filepath.Join(targetPath, "internal", "appconfig", "normalize_test.go"),
	)
	require.NoError(t, err)
	assert.Contains(t, string(normalizeTest), `t.Setenv("MYCOMMAND_TEST_HOME", "/tmp/home")`)

	readme, err := os.ReadFile(filepath.Join(targetPath, "README.md"))
	require.NoError(t, err)
	assert.Contains(t, string(readme), "# mycommand")
	assert.Contains(t, string(readme), "CLI tool that does some cool stuff")
	assert.Contains(
		t,
		string(readme),
		"generate a config file with the current\n"+
			"defaults",
	)
	assert.Contains(t, string(readme), "`~`, and `~user`")
	assert.Contains(t, string(readme), "`mycommand completion zsh`")

	exampleCmd, err := os.ReadFile(
		filepath.Join(targetPath, "cmd", "mycommand", "example", "cmd.go"),
	)
	require.NoError(t, err)
	assert.Contains(t, string(exampleCmd), "\"github.com/blumsicle/mycommand/cmd\"")
	assert.Contains(t, string(exampleCmd), "ApplyExampleOverrides")
	assert.Contains(t, string(exampleCmd), "RootPath")
	assert.Contains(t, string(exampleCmd), "GitLocation")

	completionCmd, err := os.ReadFile(
		filepath.Join(targetPath, "cmd", "mycommand", "completion", "cmd.go"),
	)
	require.NoError(t, err)
	assert.Contains(t, string(completionCmd), `enum:"zsh,bash,fish"`)
	assert.Contains(t, string(completionCmd), "var completer king.Completer")
}

func TestGenerateUsesProjectDirPrefixOnlyForTargetPath(t *testing.T) {
	rootPath := t.TempDir()

	targetPath, err := New(zerolog.Nop()).Generate(context.Background(), Config{
		Name:             "mycommand",
		Description:      "CLI tool that does some cool stuff",
		GitLocation:      "github.com/blumsicle",
		ProjectDirPrefix: "generated-",
		RootPath:         rootPath,
	})
	require.NoError(t, err)

	require.Equal(t, filepath.Join(rootPath, "generated-mycommand"), targetPath)

	goMod, err := os.ReadFile(filepath.Join(targetPath, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(goMod), "module github.com/blumsicle/mycommand")

	mainGo, err := os.ReadFile(filepath.Join(targetPath, "cmd", "mycommand", "main.go"))
	require.NoError(t, err)
	assert.Contains(t, string(mainGo), `const name = "mycommand"`)
}

func TestGenerateUsesProjectNameAsModulePathWhenGitLocationIsEmpty(t *testing.T) {
	rootPath := t.TempDir()

	targetPath, err := New(zerolog.Nop()).Generate(context.Background(), Config{
		Name:        "mycommand",
		Description: "CLI tool that does some cool stuff",
		RootPath:    rootPath,
	})
	require.NoError(t, err)

	goMod, err := os.ReadFile(filepath.Join(targetPath, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(goMod), "module mycommand")
	assert.NotContains(t, string(goMod), "module /mycommand")

	exampleCmd, err := os.ReadFile(
		filepath.Join(targetPath, "cmd", "mycommand", "example", "cmd.go"),
	)
	require.NoError(t, err)
	assert.Contains(t, string(exampleCmd), "\"mycommand/cmd\"")
	assert.Contains(t, string(exampleCmd), "ApplyExampleOverrides")
}

func TestGenerateUsesCurrentDirectoryWhenRootPathIsEmpty(t *testing.T) {
	workingDir := t.TempDir()
	originalWorkingDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(workingDir))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(originalWorkingDir))
	})

	targetPath, err := New(zerolog.Nop()).Generate(context.Background(), Config{
		Name:        "mycommand",
		Description: "CLI tool that does some cool stuff",
	})
	require.NoError(t, err)

	require.Equal(t, "mycommand", targetPath)
	assert.FileExists(t, filepath.Join(workingDir, "mycommand", "go.mod"))
}

func TestGenerateFailsWhenNameIsEmpty(t *testing.T) {
	_, err := New(zerolog.Nop()).Generate(context.Background(), Config{
		Description: "CLI tool that does some cool stuff",
		RootPath:    t.TempDir(),
	})

	require.Error(t, err)
	assert.EqualError(t, err, "name is required")
}

func TestGenerateFailsWhenDescriptionIsEmpty(t *testing.T) {
	_, err := New(zerolog.Nop()).Generate(context.Background(), Config{
		Name:     "mycommand",
		RootPath: t.TempDir(),
	})

	require.Error(t, err)
	assert.EqualError(t, err, "description is required")
}

func TestGenerateFailsWhenTargetExists(t *testing.T) {
	rootPath := t.TempDir()
	targetPath := filepath.Join(rootPath, "mycommand")
	require.NoError(t, os.MkdirAll(targetPath, 0o755))

	_, err := New(zerolog.Nop()).Generate(context.Background(), Config{
		Name:        "mycommand",
		Description: "CLI tool that does some cool stuff",
		GitLocation: "github.com/blumsicle",
		RootPath:    rootPath,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "target path already exists")
}

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

type recordingPostStep struct {
	name    string
	ran     *bool
	input   *poststep.PostStepInput
	runErr  error
	visited *[]string
}

func (s recordingPostStep) Name() string {
	return s.name
}

func (s recordingPostStep) Run(_ context.Context, input poststep.PostStepInput) error {
	if s.ran != nil {
		*s.ran = true
	}
	if s.input != nil {
		*s.input = input
	}
	if s.visited != nil {
		*s.visited = append(*s.visited, s.name)
	}
	return s.runErr
}

func TestGenerateRunsRegisteredPostSteps(t *testing.T) {
	rootPath := t.TempDir()
	gen := New(zerolog.Nop())
	var ran bool
	var input poststep.PostStepInput
	var visited []string
	gen.AddPostStep(recordingPostStep{
		name:    "record",
		ran:     &ran,
		input:   &input,
		visited: &visited,
	})

	targetPath, err := gen.Generate(context.Background(), Config{
		Name:        "mycommand",
		Description: "CLI tool that does some cool stuff",
		GitLocation: "github.com/blumsicle",
		RootPath:    rootPath,
	})
	require.NoError(t, err)

	assert.True(t, ran)
	assert.Equal(t, []string{"record"}, visited)
	assert.Equal(t, targetPath, input.ProjectPath)
	assert.Equal(t, "mycommand", input.Name)
	assert.Equal(t, "github.com/blumsicle/mycommand", input.ModulePath)
}

func TestGenerateStopsOnPostStepError(t *testing.T) {
	rootPath := t.TempDir()
	gen := New(zerolog.Nop())
	var visited []string
	gen.AddPostStep(recordingPostStep{
		name:    "first",
		runErr:  assert.AnError,
		visited: &visited,
	})
	gen.AddPostStep(recordingPostStep{
		name:    "second",
		visited: &visited,
	})

	_, err := gen.Generate(context.Background(), Config{
		Name:        "mycommand",
		Description: "CLI tool that does some cool stuff",
		RootPath:    rootPath,
	})
	require.Error(t, err)
	assert.ErrorContains(t, err, `run post step "first"`)
	assert.Equal(t, []string{"first"}, visited)
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
