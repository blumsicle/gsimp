package main

import (
	"github.com/alecthomas/kong"
	"github.com/blumsicle/gsimp/cmd"
	"github.com/blumsicle/gsimp/cmd/gsimp/create"
	"github.com/rs/zerolog"
)

type CLI struct {
	cmd.Globals

	LogLevel zerolog.Level    `short:"l" default:"info" help:"Log level"`
	Version  kong.VersionFlag `short:"v"                help:"Output version"`

	Create create.Cmd `cmd:"" help:"Create a new credentials store"`
}

func (c *CLI) GetLogLevel() zerolog.Level {
	return c.LogLevel
}

func (c *CLI) GetGlobals() any {
	return &c.Globals
}
