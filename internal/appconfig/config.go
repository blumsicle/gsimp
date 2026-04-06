// Package appconfig defines application configuration defaults and schema.
package appconfig

import "github.com/rs/zerolog"

// PostStepsConfig contains configurable post-generation step toggles.
type PostStepsConfig struct {
	GoGetUpdate bool `yaml:"go_get_update"`
	GoModTidy   bool `yaml:"go_mod_tidy"`
	GitInit     bool `yaml:"git_init"`
	GitCommit   bool `yaml:"git_commit"`
}

// Config contains app-wide settings loaded from defaults, YAML, and CLI overrides.
type Config struct {
	RootPath    string          `yaml:"root_path"`
	GitLocation string          `yaml:"git_location"`
	LogLevel    zerolog.Level   `yaml:"log_level"`
	PostSteps   PostStepsConfig `yaml:"post_steps"`
}

// Default returns a config initialized with built-in defaults.
func Default() *Config {
	return &Config{
		RootPath:    ".",
		GitLocation: "",
		LogLevel:    zerolog.InfoLevel,
		PostSteps: PostStepsConfig{
			GoGetUpdate: true,
			GoModTidy:   true,
			GitInit:     true,
			GitCommit:   true,
		},
	}
}
