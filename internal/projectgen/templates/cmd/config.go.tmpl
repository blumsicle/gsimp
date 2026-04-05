package cmd

import (
	"errors"
	"io/fs"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
)

type Config struct {
	RootPath    string        `yaml:"root_path"`
	GitLocation string        `yaml:"git_location"`
	LogLevel    zerolog.Level `yaml:"log_level"`
}

func DefaultConfig() *Config {
	return &Config{
		RootPath:    ".",
		GitLocation: "",
		LogLevel:    zerolog.InfoLevel,
	}
}

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
