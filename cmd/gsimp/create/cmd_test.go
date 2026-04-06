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

func setTestGitIdentity(t *testing.T) {
	t.Helper()
	t.Setenv("GIT_AUTHOR_NAME", "gsimp test")
	t.Setenv("GIT_AUTHOR_EMAIL", "gsimp-test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "gsimp test")
	t.Setenv("GIT_COMMITTER_EMAIL", "gsimp-test@example.com")
}

func TestRunGeneratesProject(t *testing.T) {
	setTestGitIdentity(t)

	rootPath := t.TempDir()
	gitLocation := "github.com/blumsicle"
	command := Command{
		RootPath:    &rootPath,
		GitLocation: &gitLocation,
		Name:        "cooltool",
		Description: "CLI tool that does some cool stuff",
	}

	cfg := appconfig.Default()
	cfg.RootPath = rootPath
	cfg.GitLocation = gitLocation

	err := command.Run(zerolog.Nop(), cfg)
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

func TestRunSkipsGitPostStepsWhenGitInitIsDisabled(t *testing.T) {
	setTestGitIdentity(t)

	rootPath := t.TempDir()
	gitLocation := "github.com/blumsicle"
	command := Command{
		Name:        "cooltool",
		Description: "CLI tool that does some cool stuff",
	}

	cfg := appconfig.Default()
	cfg.RootPath = rootPath
	cfg.GitLocation = gitLocation
	cfg.PostSteps.GitInit = false
	cfg.PostSteps.GitCommit = true

	err := command.Run(zerolog.Nop(), cfg)
	require.NoError(t, err)

	projectPath := filepath.Join(rootPath, "cooltool")
	assert.DirExists(t, projectPath)
	assert.NoDirExists(t, filepath.Join(projectPath, ".git"))
}

func TestRunSkipsInitialCommitWhenGitCommitIsDisabled(t *testing.T) {
	setTestGitIdentity(t)

	rootPath := t.TempDir()
	gitLocation := "github.com/blumsicle"
	command := Command{
		Name:        "cooltool",
		Description: "CLI tool that does some cool stuff",
	}

	cfg := appconfig.Default()
	cfg.RootPath = rootPath
	cfg.GitLocation = gitLocation
	cfg.PostSteps.GitCommit = false

	err := command.Run(zerolog.Nop(), cfg)
	require.NoError(t, err)

	projectPath := filepath.Join(rootPath, "cooltool")
	assert.DirExists(t, filepath.Join(projectPath, ".git"))

	commitCheck := exec.Command("git", "-C", projectPath, "rev-parse", "--verify", "HEAD")
	output, err := commitCheck.CombinedOutput()
	require.Error(t, err)
	assert.Contains(t, string(output), "Needed a single revision")
}

func TestAfterApplyOverridesConfig(t *testing.T) {
	rootPath := "/tmp/src"
	gitLocation := "github.com/acme"
	cfg := appconfig.Default()

	command := Command{
		RootPath:      &rootPath,
		GitLocation:   &gitLocation,
		NoGoGetUpdate: true,
		NoGitInit:     true,
	}

	err := command.AfterApply(cfg)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/src", cfg.RootPath)
	assert.Equal(t, "github.com/acme", cfg.GitLocation)
	assert.False(t, cfg.PostSteps.GoGetUpdate)
	assert.True(t, cfg.PostSteps.GoModTidy)
	assert.False(t, cfg.PostSteps.GitInit)
	assert.True(t, cfg.PostSteps.GitCommit)
}
