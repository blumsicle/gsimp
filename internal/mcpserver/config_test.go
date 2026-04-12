package mcpserver

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadYAMLMissingFileKeepsDefaults(t *testing.T) {
	cfg := DefaultConfig()
	err := cfg.LoadYAML(filepath.Join(t.TempDir(), "missing.yaml"))
	require.NoError(t, err)

	assert.Equal(t, "bcli", cfg.BCLICommand)
	assert.Equal(t, defaultTimeout, cfg.Timeout)
	assert.Equal(t, "info", cfg.LogLevel.String())
}

func TestLoadYAMLOverridesDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bcli-mcp.yaml")
	require.NoError(
		t,
		os.WriteFile(
			path,
			[]byte("bcli_command: /tmp/bcli\ntimeout: 2m\nlog_level: debug\n"),
			0o644,
		),
	)

	cfg := DefaultConfig()
	err := cfg.LoadYAML(path)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/bcli", cfg.BCLICommand)
	assert.Equal(t, 2*time.Minute, cfg.Timeout)
	assert.Equal(t, "debug", cfg.LogLevel.String())
}
