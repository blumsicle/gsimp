// Package main defines the bcli command-line interface.
package main

import (
	"github.com/blumsicle/bcli/cmd"
	"github.com/blumsicle/bcli/cmd/bcli/completion"
	"github.com/blumsicle/bcli/cmd/bcli/config"
	"github.com/blumsicle/bcli/cmd/bcli/create"
	"github.com/blumsicle/bcli/internal/appconfig"
)

// CLI defines the root command tree for the bcli generator.
type CLI struct {
	cmd.Globals

	Completion completion.Command `cmd:"" help:"Generate shell completions"`
	Config     config.Command     `cmd:"" help:"Write the resolved config as YAML"`
	Create     create.Command     `cmd:"" help:"Create a new Go CLI starter project"`
}

// AfterApply loads file-backed config and applies root flag overrides.
func (c *CLI) AfterApply(cfg *appconfig.Config) error {
	if err := cfg.LoadYAML(c.ConfigFile); err != nil {
		return err
	}

	cfg.ApplyRootOverrides(appconfig.RootOverrides{
		LogLevel: c.LogLevel,
	})

	return nil
}
