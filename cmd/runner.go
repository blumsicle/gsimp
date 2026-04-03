package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
)

type CLIRunner interface {
	GetLogLevel() zerolog.Level
	RunArgs() []any
}

type Config struct {
	Description string
	BuildInfo   BuildInfo
}

func Parse(cli any, cfg Config) *kong.Context {
	return kong.Parse(
		cli,
		kong.Name(cfg.BuildInfo.Name),
		kong.Description(cfg.Description),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
		kong.Vars{
			"version": fmt.Sprintf("%s %s %s", cfg.BuildInfo.Name, cfg.BuildInfo.Version, cfg.BuildInfo.Commit),
		},
	)
}

func NewLogger(level zerolog.Level) zerolog.Logger {
	cw := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime + " MST"}

	return zerolog.New(cw).
		Level(level).
		With().
		Timestamp().
		Str("logger", "main").
		Logger()
}

func Run(ctx *kong.Context, log zerolog.Logger, args ...any) error {
	runArgs := append([]any{log}, args...)
	return ctx.Run(runArgs...)
}
