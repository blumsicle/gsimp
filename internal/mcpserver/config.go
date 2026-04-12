// Package mcpserver exposes bcli project generation through MCP.
package mcpserver

import (
	"time"

	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/rs/zerolog"
)

const defaultTimeout = 10 * time.Minute

// Config defines runtime settings for the bcli MCP server.
type Config struct {
	BCLICommand string        `yaml:"bcli_command"`
	Timeout     time.Duration `yaml:"timeout"`
	LogLevel    zerolog.Level `yaml:"log_level"`
}

// DefaultConfig returns the default MCP server config.
func DefaultConfig() Config {
	return Config{
		BCLICommand: "bcli",
		Timeout:     defaultTimeout,
		LogLevel:    zerolog.InfoLevel,
	}
}

// LoadYAML loads MCP server config values from a YAML file.
func (c *Config) LoadYAML(path string) error {
	return cliutil.LoadYAML(path, c, "mcp")
}
