// Package cli contains shared runtime helpers for CLI binaries.
package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
)

// Config defines shared runtime settings for parser construction.
type Config struct {
	Description string
	BuildInfo   BuildInfo
}

// Options returns the default Kong options for a CLI runtime config.
func Options(cfg Config) []kong.Option {
	return []kong.Option{
		kong.Name(cfg.BuildInfo.Name),
		kong.Description(cfg.Description),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
		kong.Vars{
			"version": fmt.Sprintf(
				"%s %s %s",
				cfg.BuildInfo.Name,
				cfg.BuildInfo.Version,
				cfg.BuildInfo.Commit,
			),
			"command_name": cfg.BuildInfo.Name,
		},
	}
}

// New constructs a Kong parser for app with the provided runtime config.
func New(app any, cfg Config, options ...kong.Option) (*kong.Kong, error) {
	options = append(Options(cfg), options...)
	return kong.New(app, options...)
}

// Parse parses process arguments using the shared CLI runtime config.
func Parse(app any, cfg Config, options ...kong.Option) *kong.Context {
	parser, err := New(app, cfg, options...)
	if err != nil {
		panic(err)
	}

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		parser.FatalIfErrorf(err)
	}

	return ctx
}

// NewLogger constructs the standard console logger for a CLI binary.
func NewLogger(level zerolog.Level, stderr io.Writer) zerolog.Logger {
	cw := newConsoleWriter(stderr)

	return SubLogger(
		zerolog.New(cw).
			Level(level).
			With().
			Timestamp().
			Logger(),
		"main",
	)
}

func newConsoleWriter(stderr io.Writer) zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{
		Out:        stderr,
		NoColor:    !isInteractiveTerminal(stderr),
		TimeFormat: time.DateTime + " MST",
	}
}

func isInteractiveTerminal(stderr io.Writer) bool {
	file, ok := stderr.(*os.File)
	return ok && isatty.IsTerminal(file.Fd())
}

// SubLogger returns a logger rebound to the provided subsystem name.
func SubLogger(log zerolog.Logger, subsystem string) zerolog.Logger {
	return log.With().Str("logger", subsystem).Logger()
}

// Run executes the selected Kong command with the shared logger injected first.
func Run(ctx *kong.Context, log zerolog.Logger, args ...any) error {
	runArgs := append([]any{log}, args...)
	return ctx.Run(runArgs...)
}
