// Package appconfig defines application configuration defaults and schema.
package appconfig

import (
	"errors"
	"io/fs"
	"os"

	"github.com/goccy/go-yaml"
)

// LoadYAML merges YAML configuration from path into the config.
func (c *Config) LoadYAML(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	c.RootPath = os.ExpandEnv(c.RootPath)
	c.GitLocation = os.ExpandEnv(c.GitLocation)

	return nil
}
