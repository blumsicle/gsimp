package appconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	assert.Equal(t, ".", cfg.RootPath)
	assert.Equal(t, "", cfg.GitLocation)
	assert.Equal(t, zerolog.InfoLevel, cfg.LogLevel)
}

func TestLoadYAMLMissingFileKeepsDefaults(t *testing.T) {
	cfg := Default()
	err := cfg.LoadYAML(filepath.Join(t.TempDir(), "missing.yaml"))
	require.NoError(t, err)

	assert.Equal(t, ".", cfg.RootPath)
	assert.Equal(t, "", cfg.GitLocation)
	assert.Equal(t, zerolog.InfoLevel, cfg.LogLevel)
}

func TestLoadYAMLOverridesDefaults(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(
		t,
		os.WriteFile(
			configPath,
			[]byte("root_path: /tmp/src\ngit_location: github.com/acme\nlog_level: debug\n"),
			0o644,
		),
	)

	cfg := Default()
	err := cfg.LoadYAML(configPath)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/src", cfg.RootPath)
	assert.Equal(t, "github.com/acme", cfg.GitLocation)
	assert.Equal(t, zerolog.DebugLevel, cfg.LogLevel)
}

func TestLoadYAMLExpandsEnvironmentVariables(t *testing.T) {
	t.Setenv("GSIMP_TEST_HOME", "/tmp/home")

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(
		t,
		os.WriteFile(
			configPath,
			[]byte("root_path: $GSIMP_TEST_HOME/src\ngit_location: github.com/acme\n"),
			0o644,
		),
	)

	cfg := Default()
	err := cfg.LoadYAML(configPath)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/home/src", cfg.RootPath)
	assert.Equal(t, "github.com/acme", cfg.GitLocation)
}
