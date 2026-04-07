// Package main defines the bcli command-line interface.
package main

import (
	"github.com/alecthomas/kong"
	"github.com/blumsicle/bcli/cmd"
	"github.com/blumsicle/bcli/cmd/bcli/create"
	"github.com/blumsicle/bcli/internal/appconfig"
	"github.com/rs/zerolog"
)

// CLI defines the root command tree for the bcli generator.
type CLI struct {
	cmd.Globals

	LogLevel *zerolog.Level   `short:"l" help:"Log level"`
	Version  kong.VersionFlag `short:"v" help:"Output version"`

	Create create.Command `cmd:"" help:"Create a new Go CLI starter project"`
}

// AfterApply loads file-backed config and applies root flag overrides.
func (c *CLI) AfterApply(cfg *appconfig.Config) error {
	if err := cfg.LoadYAML(c.ConfigFile); err != nil {
		return err
	}

	if c.LogLevel != nil {
		cfg.LogLevel = *c.LogLevel
	}

	return nil
}
