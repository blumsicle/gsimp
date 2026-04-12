package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/blumsicle/bcli/internal/appconfig"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFileLoadsAndFlagsOverrideIt(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "bcli.yaml")
	require.NoError(
		t,
		os.WriteFile(
			configPath,
			[]byte(
				"root_path: /from-config\nproject_dir_prefix: from-config-\ngit_location: github.com/from-config\nlog_level: debug\n",
			),
			0o644,
		),
	)

	appConfig := &appconfig.Config{}
	harness := newCLITestHarness(t, appConfig)

	flagRootPath := "/from-flag"
	flagProjectDirPrefix := "from-flag-"
	flagGitLocation := "github.com/from-flag"
	flagLogLevel := zerolog.WarnLevel
	_, err := harness.parse(
		t,
		"--config-file", configPath,
		"--log-level", "warn",
		"create",
		"--root-path", flagRootPath,
		"--project-dir-prefix", flagProjectDirPrefix,
		"--git-location", flagGitLocation,
		"cooltool",
		"CLI tool that does some cool stuff",
	)
	require.NoError(t, err)

	assert.Equal(t, "/from-flag", appConfig.RootPath)
	assert.Equal(t, "from-flag-", appConfig.ProjectDirPrefix)
	assert.Equal(t, "github.com/from-flag", appConfig.GitLocation)
	assert.Equal(t, flagLogLevel, appConfig.LogLevel)
	assert.Equal(t, -1, harness.exitCode)
	assert.Empty(t, harness.stdoutString())
	assert.Empty(t, harness.stderrString())
}

func TestConfigCommandWritesMergedConfigToFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "bcli.yaml")
	outputPath := filepath.Join(t.TempDir(), "resolved.yaml")
	require.NoError(
		t,
		os.WriteFile(
			configPath,
			[]byte(
				"root_path: $HOME/from-config\ngit_location: $GIT_HOST/from-config\npost_steps:\n  git_commit: false\n",
			),
			0o644,
		),
	)

	appConfig := appconfig.Default()
	harness := newCLITestHarness(t, appConfig)

	ctx, err := harness.parse(
		t,
		"--config-file", configPath,
		"config",
		"--output", outputPath,
	)
	require.NoError(t, err)

	err = harness.run(t, ctx)
	require.NoError(t, err)

	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	var got appconfig.Config
	require.NoError(t, yaml.Unmarshal(data, &got))
	assert.Equal(t, "$HOME/from-config", got.RootPath)
	assert.Equal(t, "$GIT_HOST/from-config", got.GitLocation)
	assert.Equal(t, zerolog.InfoLevel, got.LogLevel)
	assert.True(t, got.PostSteps.GoGetUpdate)
	assert.True(t, got.PostSteps.GoModTidy)
	assert.True(t, got.PostSteps.GitInit)
	assert.False(t, got.PostSteps.GitCommit)
	assert.Equal(t, -1, harness.exitCode)
	assert.Empty(t, harness.stderrString())
}
