package main

import (
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/blumsicle/bcli/internal/bcliconfig"
	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/rs/zerolog"
)

const name = "bcli"

func main() {
	appConfig := bcliconfig.Default()
	cli := &CLI{}
	cfg := cliutil.Config{
		Description: "Generate starter Go CLI projects",
		BuildInfo:   cliutil.ResolveBuildInfo(name),
	}

	ctx := cliutil.Parse(cli, cfg, kong.Bind(&cli.Globals), kong.Bind(appConfig))

	zerolog.DurationFieldUnit = time.Minute
	zerolog.TimeFieldFormat = time.DateTime + " MST"

	log := cliutil.NewLogger(appConfig.LogLevel, ctx.Stderr)

	if err := cliutil.Run(ctx, log); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
