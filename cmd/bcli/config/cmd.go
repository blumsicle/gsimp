// Package config implements the bcli config subcommand.
package config

import (
	"fmt"
	"os"

	"github.com/blumsicle/bcli/internal/appconfig"
	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
)

// Command writes the resolved application config as YAML.
type Command struct {
	Output string `short:"o" default:"-" type:"path" help:"Path to write the resolved config YAML, or - for stdout"`
}

// Run writes the merged config to stdout or a file.
func (c *Command) Run(log zerolog.Logger, cfg *appconfig.Config) (err error) {
	log = cliutil.SubLogger(log, "config")

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config yaml: %w", err)
	}

	outputPath := c.Output
	if outputPath == "" {
		outputPath = "-"
	}

	log.Debug().
		Str("output", outputPath).
		Bool("stdout", outputPath == "-").
		Str("root_path", cfg.RootPath).
		Str("git_location", cfg.GitLocation).
		Msg("resolved config command configuration")

	output := os.Stdout
	if c.Output != "" && c.Output != "-" {
		file, err := os.Create(c.Output)
		if err != nil {
			return fmt.Errorf("create config file: %w", err)
		}
		defer func() {
			closeErr := file.Close()
			if err == nil && closeErr != nil {
				err = fmt.Errorf("close config file: %w", closeErr)
			}
		}()
		output = file
	}

	if _, err := output.Write(data); err != nil {
		return fmt.Errorf("write config yaml: %w", err)
	}

	log.Info().
		Str("output", outputPath).
		Msg("wrote resolved config")

	return nil
}
