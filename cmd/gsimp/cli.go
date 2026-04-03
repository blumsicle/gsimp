package main

import (
	"github.com/alecthomas/kong"
	"github.com/blumsicle/gsimp/cmd"
	"github.com/blumsicle/gsimp/cmd/gsimp/example"
	"github.com/rs/zerolog"
)

type CLI struct {
	cmd.Globals

	LogLevel zerolog.Level    `short:"l" default:"info" help:"Log level"`
	Version  kong.VersionFlag `short:"v"                help:"Output version"`

	Example example.Command `cmd:"" help:"Example subcommand for new projects"`
}

func (c *CLI) GetLogLevel() zerolog.Level {
	return c.LogLevel
}

func (c *CLI) RunArgs() []any {
	return []any{&c.Globals}
}
