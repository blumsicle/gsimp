// Package cmd contains shared command-line types that are imported by bcli subcommands.
package cmd

// Globals defines root-level flags that are injected into command handlers.
type Globals struct {
	ConfigFile string `short:"c" default:"~/.config/bcli/config.yaml" type:"path" help:"Path to the config file"`
}
