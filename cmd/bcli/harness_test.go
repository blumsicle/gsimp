package main

import (
	"bytes"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/blumsicle/bcli/internal/bcliconfig"
	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/rs/zerolog"
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
	appConfig *bcliconfig.Config
	stdout    bytes.Buffer
	stderr    bytes.Buffer
	exitCode  int
	parser    *kong.Kong
}

func newCLITestHarness(t *testing.T, appConfig *bcliconfig.Config) *cliTestHarness {
	t.Helper()
	if appConfig == nil {
		appConfig = &bcliconfig.Config{}
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
