package main

import (
	"os"

	"github.com/blumsicle/gsimp/cmd"
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
	log := cmd.NewLogger(cli.GetLogLevel())

	if err := cmd.Run(ctx, log, cli.RunArgs()...); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
