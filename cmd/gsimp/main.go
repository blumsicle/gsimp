package main

import (
	"os"
	"time"

	cliutil "github.com/blumsicle/gsimp/internal/cli"
	"github.com/rs/zerolog"
)

var (
	name    = "gsimp"
	version = "dev"
	commit  = "unknown"
)

func main() {
	cli := &CLI{}
	cfg := cliutil.Config{
		Description: "Starter CLI template",
		BuildInfo: cliutil.BuildInfo{
			Name:    name,
			Version: version,
			Commit:  commit,
		},
	}

	ctx := cliutil.Parse(cli, cfg)

	zerolog.DurationFieldUnit = time.Minute
	zerolog.TimeFieldFormat = time.DateTime + " MST"

	log := cliutil.NewLogger(cli.GetLogLevel())

	if err := cliutil.Run(ctx, log, cli.RunArgs()...); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
