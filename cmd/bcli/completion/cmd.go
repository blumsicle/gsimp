// Package completion implements shell completion commands for bcli.
package completion

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/miekg/king"
)

// Command prints shell completion definitions for supported shells.
type Command struct {
	Shell string `arg:"" enum:"zsh,bash,fish" help:"Shell to generate completions for"`
}

// Run writes the requested shell completion script to stdout.
func (c *Command) Run(ctx *kong.Context) error {
	var completer king.Completer

	switch c.Shell {
	case "zsh":
		completer = &king.Zsh{}
	case "bash":
		completer = &king.Bash{}
	case "fish":
		completer = &king.Fish{}
	default:
		return fmt.Errorf("unsupported shell %q", c.Shell)
	}

	completer.Completion(ctx.Model.Node, ctx.Model.Name)
	if _, err := ctx.Stdout.Write(completer.Out()); err != nil {
		return fmt.Errorf("write %s completion: %w", c.Shell, err)
	}
	return nil
}
