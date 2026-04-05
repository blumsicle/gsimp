package create

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/blumsicle/gsimp/internal/appconfig"
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

	err := command.Run(zerolog.Nop(), &appconfig.Config{
		RootPath:    rootPath,
		GitLocation: gitLocation,
	})
	require.NoError(t, err)

	projectPath := filepath.Join(rootPath, "cooltool")
	assert.DirExists(t, projectPath)
	assert.DirExists(t, filepath.Join(projectPath, ".git"))
	assert.FileExists(t, filepath.Join(projectPath, "go.mod"))
	assert.FileExists(t, filepath.Join(projectPath, "go.sum"))
	assert.FileExists(t, filepath.Join(projectPath, "Makefile"))
	assert.FileExists(t, filepath.Join(projectPath, "cmd", "cooltool", "main.go"))

	gitConfig, err := os.ReadFile(filepath.Join(projectPath, ".git", "config"))
	require.NoError(t, err)
	assert.Contains(t, string(gitConfig), "[core]")

	commitMessage, err := exec.Command("git", "-C", projectPath, "log", "-1", "--pretty=%s").
		CombinedOutput()
	require.NoError(t, err, string(commitMessage))
	assert.Equal(t, "Initial commit\n", string(commitMessage))
}

func TestAfterApplyOverridesConfig(t *testing.T) {
	rootPath := "/tmp/src"
	gitLocation := "github.com/acme"
	cfg := appconfig.Default()

	command := Command{
		RootPath:    &rootPath,
		GitLocation: &gitLocation,
	}

	err := command.AfterApply(cfg)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/src", cfg.RootPath)
	assert.Equal(t, "github.com/acme", cfg.GitLocation)
}
