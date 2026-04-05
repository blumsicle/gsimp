package projectgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCreatesStarterProject(t *testing.T) {
	rootPath := t.TempDir()

	targetPath, err := New().Generate(Config{
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
	assert.FileExists(t, filepath.Join(targetPath, "cmd", "config.go"))
	assert.FileExists(t, filepath.Join(targetPath, "cmd", "config_test.go"))

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

	targetPath, err := New().Generate(Config{
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

	_, err := New().Generate(Config{
		Name:        "mycommand",
		Description: "CLI tool that does some cool stuff",
		GitLocation: "github.com/blumsicle",
		RootPath:    rootPath,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "target path already exists")
}
