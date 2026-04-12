package main

import (
	"os"
	"time"

	"github.com/alecthomas/kong"
	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/blumsicle/bcli/internal/mcpserver"
	"github.com/rs/zerolog"
)

const name = "bcli-mcp"

func main() {
	appConfig := mcpserver.DefaultConfig()
	cli := &CLI{}
	cfg := cliutil.Config{
		Description: "Run the bcli MCP server",
		BuildInfo:   cliutil.ResolveBuildInfo(name),
	}

	ctx := cliutil.Parse(cli, cfg, kong.Bind(&cli.Globals), kong.Bind(&appConfig))

	zerolog.DurationFieldUnit = time.Minute
	zerolog.TimeFieldFormat = time.DateTime + " MST"

	log := cliutil.NewLogger(appConfig.LogLevel)

	if err := cliutil.Run(ctx, log); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
