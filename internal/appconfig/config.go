// Package appconfig defines application configuration defaults and schema.
package appconfig

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

// PostStepsConfig contains configurable post-generation step toggles.
type PostStepsConfig struct {
	GoGetUpdate bool `yaml:"go_get_update"`
	GoModTidy   bool `yaml:"go_mod_tidy"`
	GitInit     bool `yaml:"git_init"`
	GitCommit   bool `yaml:"git_commit"`
}

// Config contains app-wide settings loaded from defaults, YAML, and CLI overrides.
type Config struct {
	RootPath         string          `yaml:"root_path"`
	ProjectDirPrefix string          `yaml:"project_dir_prefix"`
	GitLocation      string          `yaml:"git_location"`
	LogLevel         zerolog.Level   `yaml:"log_level"`
	PostSteps        PostStepsConfig `yaml:"post_steps"`
}

// Default returns a config initialized with built-in defaults.
func Default() *Config {
	return &Config{
		RootPath:         ".",
		ProjectDirPrefix: "",
		GitLocation:      "",
		LogLevel:         zerolog.InfoLevel,
		PostSteps: PostStepsConfig{
			GoGetUpdate: true,
			GoModTidy:   true,
			GitInit:     true,
			GitCommit:   true,
		},
	}
}

// Normalize expands config values that should be resolved before command execution.
func (c *Config) Normalize() {
	c.RootPath = expandConfigPath(c.RootPath)
}

func expandConfigPath(value string) string {
	value = os.ExpandEnv(value)
	tilde, rest, hasRest := strings.Cut(value, "/")
	if !strings.HasPrefix(tilde, "~") {
		return value
	}

	homeDir, ok := tildeHomeDir(tilde)
	if !ok {
		return value
	}
	if !hasRest {
		return homeDir
	}

	return filepath.Join(homeDir, rest)
}

func tildeHomeDir(tilde string) (string, bool) {
	if tilde == "~" {
		homeDir, err := os.UserHomeDir()
		return homeDir, err == nil
	}

	name := strings.TrimPrefix(tilde, "~")
	u, err := user.Lookup(name)
	if err != nil {
		return "", false
	}

	return u.HomeDir, true
}
