package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/blumsicle/bcli/internal/appconfig"
	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConfig() cliutil.Config {
	return cliutil.Config{
		Description: "Generate starter Go CLI projects",
		BuildInfo: cliutil.BuildInfo{
			Name:    "bcli",
			Version: "test-version",
			Commit:  "test-commit",
		},
	}
}

type cliTestHarness struct {
	cli       *CLI
	appConfig *appconfig.Config
	stdout    bytes.Buffer
	stderr    bytes.Buffer
	exitCode  int
	parser    *kong.Kong
}

func newCLITestHarness(t *testing.T, appConfig *appconfig.Config) *cliTestHarness {
	t.Helper()
	if appConfig == nil {
		appConfig = &appconfig.Config{}
	}

	harness := &cliTestHarness{
		cli:       &CLI{},
		appConfig: appConfig,
		exitCode:  -1,
	}

	parser, err := cliutil.New(
		harness.cli,
		testConfig(),
		kong.Bind(&harness.cli.Globals),
		kong.Bind(harness.appConfig),
		kong.Writers(&harness.stdout, &harness.stderr),
		kong.Exit(func(code int) {
			harness.exitCode = code
		}),
	)
	require.NoError(t, err)

	harness.parser = parser
	return harness
}

func (h *cliTestHarness) parse(t *testing.T, args ...string) (*kong.Context, error) {
	t.Helper()
	return h.parser.Parse(args)
}

func (h *cliTestHarness) run(t *testing.T, ctx *kong.Context) error {
	t.Helper()
	log := zerolog.New(&bytes.Buffer{})
	return cliutil.Run(ctx, log)
}

func (h *cliTestHarness) stdoutString() string {
	return h.stdout.String()
}

func (h *cliTestHarness) stderrString() string {
	return h.stderr.String()
}

func TestVersionFlag(t *testing.T) {
	harness := newCLITestHarness(t, nil)

	_, err := harness.parse(t, "--version")
	require.Error(t, err)

	assert.Equal(t, 0, harness.exitCode)
	assert.Equal(t, "bcli test-version test-commit\n", harness.stdoutString())
	assert.Empty(t, harness.stderrString())
}

func TestHelpFlag(t *testing.T) {
	harness := newCLITestHarness(t, nil)

	_, err := harness.parse(t, "--help")
	require.Error(t, err)

	assert.Equal(t, 0, harness.exitCode)
	assert.Contains(t, harness.stdoutString(), "Generate starter Go CLI projects")
	assert.Contains(t, harness.stdoutString(), "--log-level")
	assert.Contains(t, harness.stdoutString(), "completion")
	assert.Contains(t, harness.stdoutString(), "config")
	assert.Contains(t, harness.stdoutString(), "create")
	assert.Empty(t, harness.stderrString())
}

func TestCompletionCommandWritesShellCompletionScript(t *testing.T) {
	tests := []struct {
		shell    string
		contains []string
	}{
		{
			shell: "zsh",
			contains: []string{
				"#compdef bcli",
				"compdef _bcli bcli",
				"_bcli() {",
			},
		},
		{
			shell: "bash",
			contains: []string{
				"_bcli_completions()",
				"complete -F _bcli_completions bcli",
			},
		},
		{
			shell: "fish",
			contains: []string{
				"# fish shell completion for bcli",
				"complete -c bcli -f",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			harness := newCLITestHarness(t, nil)

			ctx, err := harness.parse(t, "completion", tt.shell)
			require.NoError(t, err)

			err = harness.run(t, ctx)
			require.NoError(t, err)

			assert.Equal(t, -1, harness.exitCode)
			for _, want := range tt.contains {
				assert.Contains(t, harness.stdoutString(), want)
			}
			assert.Empty(t, harness.stderrString())
		})
	}
}

func TestConfigFileLoadsAndFlagsOverrideIt(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
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
	configPath := filepath.Join(t.TempDir(), "config.yaml")
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
