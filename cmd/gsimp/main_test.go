package main

import (
	"bytes"
	"testing"

	"github.com/alecthomas/kong"
	cliutil "github.com/blumsicle/gsimp/internal/cli"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConfig() cliutil.Config {
	return cliutil.Config{
		Description: "Starter CLI template",
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
	stdout *bytes.Buffer,
	stderr *bytes.Buffer,
	exitCode *int,
) *kong.Kong {
	t.Helper()

	parser, err := cliutil.New(
		cli,
		testConfig(),
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := -1

	parser := newTestParser(t, cli, &stdout, &stderr, &exitCode)

	_, err := parser.Parse([]string{"--version"})
	require.Error(t, err)

	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "gsimp test-version test-commit\n", stdout.String())
	assert.Empty(t, stderr.String())
}

func TestHelpFlag(t *testing.T) {
	cli := &CLI{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := -1

	parser := newTestParser(t, cli, &stdout, &stderr, &exitCode)

	_, err := parser.Parse([]string{"--help"})
	require.Error(t, err)

	assert.Equal(t, 0, exitCode)
	assert.Contains(t, stdout.String(), "Starter CLI template")
	assert.Contains(t, stdout.String(), "--config-file")
	assert.Contains(t, stdout.String(), "example")
	assert.Empty(t, stderr.String())
}

func TestExampleCommandReceivesInjectedGlobals(t *testing.T) {
	cli := &CLI{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := -1

	parser := newTestParser(t, cli, &stdout, &stderr, &exitCode)

	ctx, err := parser.Parse([]string{"--config-file", "/tmp/test-config.yaml", "example"})
	require.NoError(t, err)

	var logs bytes.Buffer
	log := zerolog.New(&logs)

	err = cliutil.Run(ctx, log, cli.RunArgs()...)
	require.NoError(t, err)

	assert.Equal(t, -1, exitCode)
	assert.Contains(t, logs.String(), "example command")
	assert.Contains(t, logs.String(), "/tmp/test-config.yaml")
}
