// Package cli contains shared runtime helpers for CLI binaries.
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
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
func NewLogger(level zerolog.Level) zerolog.Logger {
	cw := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime + " MST"}

	return zerolog.New(cw).
		Level(level).
		With().
		Timestamp().
		Str("logger", "main").
		Logger()
}

// Run executes the selected Kong command with the shared logger injected first.
func Run(ctx *kong.Context, log zerolog.Logger, args ...any) error {
	runArgs := append([]any{log}, args...)
	return ctx.Run(runArgs...)
}
