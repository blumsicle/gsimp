package appconfig

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

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
