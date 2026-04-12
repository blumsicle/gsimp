package bcliconfig

import "github.com/rs/zerolog"

// RootOverrides contains root command flag values that can override config.
type RootOverrides struct {
	LogLevel *zerolog.Level
}

// ApplyRootOverrides applies root command flag overrides to the config.
func (c *Config) ApplyRootOverrides(overrides RootOverrides) {
	if overrides.LogLevel != nil {
		c.LogLevel = *overrides.LogLevel
	}
}
