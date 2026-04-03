package main

import (
	"os"
	"time"

	"github.com/blumsicle/gsimp/cmd"
	"github.com/rs/zerolog"
)

var (
	name    = "gsimp"
	version = "dev"
	commit  = "unknown"
)

func main() {
	cli := &CLI{}
	cfg := cmd.Config{
		Description: "Credential manager",
		BuildInfo: cmd.BuildInfo{
			Name:    name,
			Version: version,
			Commit:  commit,
		},
	}

	ctx := cmd.Parse(cli, cfg)

	zerolog.DurationFieldUnit = time.Minute
	zerolog.TimeFieldFormat = time.DateTime + " MST"

	log := cmd.NewLogger(cli.GetLogLevel())

	if err := cmd.Run(ctx, log, cli.RunArgs()...); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
