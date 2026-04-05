package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/blumsicle/gsimp/internal/appconfig"
	cliutil "github.com/blumsicle/gsimp/internal/cli"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConfig() cliutil.Config {
	return cliutil.Config{
		Description: "Generate starter Go CLI projects",
		BuildInfo: cliutil.BuildInfo{
			Name:    "gsimp",
			Version: "test-version",
			Commit:  "test-commit",
		},
	}
}

func newTestParser(
	t *testing.T,
	cli *CLI,
	appConfig *appconfig.Config,
	stdout *bytes.Buffer,
	stderr *bytes.Buffer,
	exitCode *int,
) *kong.Kong {
	t.Helper()

	parser, err := cliutil.New(
		cli,
		testConfig(),
		kong.Bind(&cli.Globals),
		kong.Bind(appConfig),
		kong.Writers(stdout, stderr),
		kong.Exit(func(code int) {
			*exitCode = code
		}),
	)
	require.NoError(t, err)

	return parser
}

func TestVersionFlag(t *testing.T) {
	cli := &CLI{}
	appConfig := &appconfig.Config{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := -1

	parser := newTestParser(t, cli, appConfig, &stdout, &stderr, &exitCode)

	_, err := parser.Parse([]string{"--version"})
	require.Error(t, err)

	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "gsimp test-version test-commit\n", stdout.String())
	assert.Empty(t, stderr.String())
}

func TestHelpFlag(t *testing.T) {
	cli := &CLI{}
	appConfig := &appconfig.Config{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := -1

	parser := newTestParser(t, cli, appConfig, &stdout, &stderr, &exitCode)

	_, err := parser.Parse([]string{"--help"})
	require.Error(t, err)

	assert.Equal(t, 0, exitCode)
	assert.Contains(t, stdout.String(), "Generate starter Go CLI projects")
	assert.Contains(t, stdout.String(), "--log-level")
	assert.Contains(t, stdout.String(), "create")
	assert.Empty(t, stderr.String())
}

func TestConfigFileLoadsAndFlagsOverrideIt(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(
		t,
		os.WriteFile(
			configPath,
			[]byte(
				"root_path: /from-config\ngit_location: github.com/from-config\nlog_level: debug\n",
			),
			0o644,
		),
	)

	cli := &CLI{}
	appConfig := &appconfig.Config{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := -1

	parser := newTestParser(t, cli, appConfig, &stdout, &stderr, &exitCode)

	flagRootPath := "/from-flag"
	flagGitLocation := "github.com/from-flag"
	flagLogLevel := zerolog.WarnLevel
	_, err := parser.Parse([]string{
		"--config-file", configPath,
		"--log-level", "warn",
		"create",
		"--root-path", flagRootPath,
		"--git-location", flagGitLocation,
		"cooltool",
		"CLI tool that does some cool stuff",
	})
	require.NoError(t, err)

	assert.Equal(t, "/from-flag", appConfig.RootPath)
	assert.Equal(t, "github.com/from-flag", appConfig.GitLocation)
	assert.Equal(t, flagLogLevel, appConfig.LogLevel)
	assert.Equal(t, -1, exitCode)
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())
}
