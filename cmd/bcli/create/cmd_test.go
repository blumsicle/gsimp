package create

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/blumsicle/bcli/internal/bcliconfig"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCLI struct {
	Create Command `cmd:""`
}

func captureStdout(t *testing.T) (*os.File, func() []byte) {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = writer

	return writer, func() []byte {
		require.NoError(t, writer.Close())
		os.Stdout = originalStdout

		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.NoError(t, reader.Close())

		return data
	}
}

func setTestGitIdentity(t *testing.T) {
	t.Helper()
	t.Setenv("GIT_AUTHOR_NAME", "bcli test")
	t.Setenv("GIT_AUTHOR_EMAIL", "bcli-test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "bcli test")
	t.Setenv("GIT_COMMITTER_EMAIL", "bcli-test@example.com")
}

func TestRunGeneratesProject(t *testing.T) {
	setTestGitIdentity(t)

	rootPath := t.TempDir()
	gitLocation := "github.com/blumsicle"
	projectDirPrefix := "generated-"
	command := Command{
		RootPath:         &rootPath,
		ProjectDirPrefix: &projectDirPrefix,
		GitLocation:      &gitLocation,
		Name:             "cooltool",
		Description:      "CLI tool that does some cool stuff",
	}

	cfg := bcliconfig.Default()
	cfg.RootPath = rootPath
	cfg.ProjectDirPrefix = projectDirPrefix
	cfg.GitLocation = gitLocation

	err := command.Run(zerolog.Nop(), cfg)
	require.NoError(t, err)

	projectPath := filepath.Join(rootPath, "generated-cooltool")
	assert.DirExists(t, projectPath)
	assert.DirExists(t, filepath.Join(projectPath, ".git"))
	assert.FileExists(t, filepath.Join(projectPath, "go.mod"))
	assert.FileExists(t, filepath.Join(projectPath, "go.sum"))
	assert.FileExists(t, filepath.Join(projectPath, "Taskfile.yml"))
	assert.FileExists(t, filepath.Join(projectPath, "cmd", "cooltool", "main.go"))

	gitConfig, err := os.ReadFile(filepath.Join(projectPath, ".git", "config"))
	require.NoError(t, err)
	assert.Contains(t, string(gitConfig), "[core]")

	commitMessage, err := exec.Command("git", "-C", projectPath, "log", "-1", "--pretty=%s").
		CombinedOutput()
	require.NoError(t, err, string(commitMessage))
	assert.Equal(t, "Initial commit\n", string(commitMessage))
}

func TestRunWritesCreateResultAsJSON(t *testing.T) {
	rootPath := t.TempDir()
	gitLocation := "github.com/blumsicle"
	command := Command{
		NoGoGetUpdate: true,
		NoGoModTidy:   true,
		NoGitInit:     true,
		NoGitCommit:   true,
		JSON:          true,
		Name:          "cooltool",
		Description:   "CLI tool that does some cool stuff",
	}

	cfg := bcliconfig.Default()
	cfg.RootPath = rootPath
	cfg.GitLocation = gitLocation
	cfg.PostSteps.GoGetUpdate = false
	cfg.PostSteps.GoModTidy = false
	cfg.PostSteps.GitInit = false
	cfg.PostSteps.GitCommit = false

	_, restoreStdout := captureStdout(t)

	err := command.Run(zerolog.Nop(), cfg)
	require.NoError(t, err)

	var got CreateResult
	require.NoError(t, json.Unmarshal(restoreStdout(), &got))
	assert.Equal(t, "cooltool", got.Project)
	assert.Equal(t, "CLI tool that does some cool stuff", got.Description)
	assert.Equal(t, "github.com/blumsicle/cooltool", got.ModulePath)
	assert.Equal(t, filepath.Join(rootPath, "cooltool"), got.TargetPath)
	assert.False(t, got.InPlace)
	assert.Equal(
		t,
		[]PostStepResult{
			{Name: "go get -u ./...", Ran: false},
			{Name: "go mod tidy", Ran: false},
			{Name: "git init", Ran: false},
			{Name: "git commit", Ran: false},
		},
		got.PostSteps,
	)
}

func TestRunCreatesProjectInPlace(t *testing.T) {
	workingDir := t.TempDir()
	originalWorkingDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(workingDir))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(originalWorkingDir))
	})

	command := Command{
		NoGoGetUpdate: true,
		NoGoModTidy:   true,
		NoGitInit:     true,
		NoGitCommit:   true,
		InPlace:       true,
		JSON:          true,
		Name:          "cooltool",
		Description:   "CLI tool that does some cool stuff",
	}

	cfg := bcliconfig.Default()
	cfg.PostSteps.GoGetUpdate = false
	cfg.PostSteps.GoModTidy = false
	cfg.PostSteps.GitInit = false
	cfg.PostSteps.GitCommit = false
	cfg.RootPath = filepath.Join(t.TempDir(), "ignored-root")
	cfg.ProjectDirPrefix = "ignored-prefix-"

	_, restoreStdout := captureStdout(t)

	err = command.Run(zerolog.Nop(), cfg)
	require.NoError(t, err)

	var got CreateResult
	require.NoError(t, json.Unmarshal(restoreStdout(), &got))
	expectedTargetPath, err := filepath.EvalSymlinks(workingDir)
	require.NoError(t, err)
	assert.Equal(t, expectedTargetPath, got.TargetPath)
	assert.True(t, got.InPlace)
	assert.FileExists(t, filepath.Join(workingDir, "go.mod"))
	assert.FileExists(t, filepath.Join(workingDir, "cmd", "cooltool", "main.go"))
	assert.NoDirExists(t, filepath.Join(workingDir, "ignored-prefix-cooltool"))
}

func TestParseRejectsInPlaceWithProjectDirPrefixFlag(t *testing.T) {
	parser, err := kong.New(&testCLI{})
	require.NoError(t, err)

	_, err = parser.Parse([]string{
		"create",
		"--inplace",
		"--project-dir-prefix",
		"generated-",
		"cooltool",
		"CLI tool that does some cool stuff",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "can't be used together")
	assert.Contains(t, err.Error(), "--inplace")
	assert.Contains(t, err.Error(), "--project-dir-prefix")
}

func TestRunSkipsGitPostStepsWhenGitInitIsDisabled(t *testing.T) {
	setTestGitIdentity(t)

	rootPath := t.TempDir()
	gitLocation := "github.com/blumsicle"
	command := Command{
		Name:        "cooltool",
		Description: "CLI tool that does some cool stuff",
	}

	cfg := bcliconfig.Default()
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

	cfg := bcliconfig.Default()
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
	projectDirPrefix := "generated-"
	gitLocation := "github.com/acme"
	cfg := bcliconfig.Default()

	command := Command{
		RootPath:         &rootPath,
		ProjectDirPrefix: &projectDirPrefix,
		GitLocation:      &gitLocation,
		NoGoGetUpdate:    true,
		NoGitInit:        true,
	}

	err := command.AfterApply(cfg)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/src", cfg.RootPath)
	assert.Equal(t, "generated-", cfg.ProjectDirPrefix)
	assert.Equal(t, "github.com/acme", cfg.GitLocation)
	assert.False(t, cfg.PostSteps.GoGetUpdate)
	assert.True(t, cfg.PostSteps.GoModTidy)
	assert.False(t, cfg.PostSteps.GitInit)
	assert.True(t, cfg.PostSteps.GitCommit)
}

func TestRunExpandsEnvironmentVariablesForGeneration(t *testing.T) {
	setTestGitIdentity(t)
	t.Setenv("BCLI_CREATE_ROOT", t.TempDir())
	t.Setenv("BCLI_GIT_HOST", "github.com")

	command := Command{
		Name:        "cooltool",
		Description: "CLI tool that does some cool stuff",
	}

	cfg := bcliconfig.Default()
	cfg.RootPath = "$BCLI_CREATE_ROOT"
	cfg.GitLocation = "$BCLI_GIT_HOST/blumsicle"
	cfg.PostSteps.GoGetUpdate = false
	cfg.PostSteps.GoModTidy = false
	cfg.PostSteps.GitInit = false
	cfg.PostSteps.GitCommit = false

	err := command.Run(zerolog.Nop(), cfg)
	require.NoError(t, err)

	projectPath := filepath.Join(os.ExpandEnv(cfg.RootPath), "cooltool")
	assert.DirExists(t, projectPath)

	goMod, err := os.ReadFile(filepath.Join(projectPath, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(goMod), "module $BCLI_GIT_HOST/blumsicle/cooltool")
}
