package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
)

type Runner interface {
	GetLogLevel() zerolog.Level
	RunArgs() []any
}

type Config struct {
	Description string
	BuildInfo   BuildInfo
}

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

func New(app any, cfg Config, options ...kong.Option) (*kong.Kong, error) {
	options = append(Options(cfg), options...)
	return kong.New(app, options...)
}

func Parse(app any, cfg Config) *kong.Context {
	parser, err := New(app, cfg)
	if err != nil {
		panic(err)
	}

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		parser.FatalIfErrorf(err)
	}

	return ctx
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
