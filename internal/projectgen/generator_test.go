package projectgen

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/blumsicle/gsimp/internal/poststep"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCreatesStarterProject(t *testing.T) {
	rootPath := t.TempDir()

	targetPath, err := New(zerolog.Nop()).Generate(context.Background(), Config{
		Name:        "mycommand",
		Description: "CLI tool that does some cool stuff",
		GitLocation: "github.com/blumsicle",
		RootPath:    rootPath,
	})
	require.NoError(t, err)

	require.Equal(t, filepath.Join(rootPath, "mycommand"), targetPath)

	goMod, err := os.ReadFile(filepath.Join(targetPath, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(goMod), "module github.com/blumsicle/mycommand")

	mainGo, err := os.ReadFile(filepath.Join(targetPath, "cmd", "mycommand", "main.go"))
	require.NoError(t, err)
	assert.Contains(t, string(mainGo), "Description: \"CLI tool that does some cool stuff\"")
	assert.Contains(t, string(mainGo), "name    = \"mycommand\"")

	assert.FileExists(t, filepath.Join(targetPath, "mycommand.yaml"))
	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "config.go"))
	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "load.go"))
	assert.FileExists(t, filepath.Join(targetPath, "internal", "appconfig", "config_test.go"))

	readme, err := os.ReadFile(filepath.Join(targetPath, "README.md"))
	require.NoError(t, err)
	assert.Contains(t, string(readme), "# mycommand")
	assert.Contains(t, string(readme), "CLI tool that does some cool stuff")

	exampleCmd, err := os.ReadFile(
		filepath.Join(targetPath, "cmd", "mycommand", "example", "cmd.go"),
	)
	require.NoError(t, err)
	assert.Contains(t, string(exampleCmd), "\"github.com/blumsicle/mycommand/cmd\"")
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
