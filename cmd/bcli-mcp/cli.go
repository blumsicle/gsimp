// Package main defines the bcli-mcp command-line interface.
package main

import (
	"context"

	"github.com/blumsicle/bcli/cmd"
	"github.com/blumsicle/bcli/internal/mcpserver"
	"github.com/rs/zerolog"
)

// CLI defines the root command for the bcli-mcp server.
type CLI struct {
	cmd.Globals
}

// AfterApply loads file-backed MCP server config and applies root flag overrides.
func (c *CLI) AfterApply(cfg *mcpserver.Config) error {
	if err := cfg.LoadYAML(c.ConfigFile); err != nil {
		return err
	}
	if c.LogLevel != nil {
		cfg.LogLevel = *c.LogLevel
	}

	return nil
}

// Run starts the MCP server over stdio.
func (c *CLI) Run(log zerolog.Logger, cfg *mcpserver.Config) error {
	return mcpserver.New(*cfg, nil).RunStdio(context.Background())
}
