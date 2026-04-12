// Package bcliconfig defines bcli configuration defaults and schema.
package bcliconfig

import cliutil "github.com/blumsicle/bcli/internal/cli"

// LoadYAML merges YAML configuration from path into the config.
func (c *Config) LoadYAML(path string) error {
	return cliutil.LoadYAML(path, c, "bcli")
}
