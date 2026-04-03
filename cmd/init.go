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
	GetGlobals() any
}

func Run(cli CLIRunner, description string, info BuildInfo) {
	ctx := kong.Parse(
		cli,
		kong.Name(info.Name),
		kong.Description(description),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
		kong.Vars{
			"version": fmt.Sprintf("%s %s %s", info.Name, info.Version, info.Commit),
		},
	)

	cw := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime + " MST"}

	zerolog.SetGlobalLevel(cli.GetLogLevel())
	zerolog.DurationFieldUnit = time.Minute
	zerolog.TimeFieldFormat = time.DateTime + " MST"

	log := zerolog.New(cw).With().Timestamp().Str("logger", "main").Logger()

	err := ctx.Run(log, cli.GetGlobals())
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
