package create

import (
	"path/filepath"
	"testing"

	"github.com/blumsicle/gsimp/cmd"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunGeneratesProject(t *testing.T) {
	rootPath := t.TempDir()
	gitLocation := "github.com/blumsicle"
	command := Command{
		RootPath:    &rootPath,
		GitLocation: &gitLocation,
		Name:        "cooltool",
		Description: "CLI tool that does some cool stuff",
	}

	err := command.Run(zerolog.Nop(), &cmd.Config{
		RootPath:    rootPath,
		GitLocation: gitLocation,
	})
	require.NoError(t, err)

	projectPath := filepath.Join(rootPath, "cooltool")
	assert.DirExists(t, projectPath)
	assert.FileExists(t, filepath.Join(projectPath, "go.mod"))
	assert.FileExists(t, filepath.Join(projectPath, "Makefile"))
	assert.FileExists(t, filepath.Join(projectPath, "cmd", "cooltool", "main.go"))
	assert.FileExists(t, filepath.Join(projectPath, "cmd", "cooltool", "example", "cmd.go"))
	assert.FileExists(t, filepath.Join(projectPath, "internal", "cli", "runner.go"))
	assert.FileExists(t, filepath.Join(projectPath, "README.md"))
}

func TestAfterApplyOverridesConfig(t *testing.T) {
	rootPath := "/tmp/src"
	gitLocation := "github.com/acme"
	cfg := cmd.DefaultConfig()

	command := Command{
		RootPath:    &rootPath,
		GitLocation: &gitLocation,
	}

	err := command.AfterApply(cfg)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/src", cfg.RootPath)
	assert.Equal(t, "github.com/acme", cfg.GitLocation)
}
