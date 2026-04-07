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
	assert.True(t, cfg.PostSteps.GoGetUpdate)
	assert.True(t, cfg.PostSteps.GoModTidy)
	assert.True(t, cfg.PostSteps.GitInit)
	assert.True(t, cfg.PostSteps.GitCommit)
}

func TestLoadYAMLMissingFileKeepsDefaults(t *testing.T) {
	cfg := Default()
	err := cfg.LoadYAML(filepath.Join(t.TempDir(), "missing.yaml"))
	require.NoError(t, err)

	assert.Equal(t, ".", cfg.RootPath)
	assert.Equal(t, "", cfg.GitLocation)
	assert.Equal(t, zerolog.InfoLevel, cfg.LogLevel)
	assert.True(t, cfg.PostSteps.GoGetUpdate)
	assert.True(t, cfg.PostSteps.GoModTidy)
	assert.True(t, cfg.PostSteps.GitInit)
	assert.True(t, cfg.PostSteps.GitCommit)
}

func TestLoadYAMLOverridesDefaults(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(
		t,
		os.WriteFile(
			configPath,
			[]byte(
				"root_path: /tmp/src\n"+
					"git_location: github.com/acme\n"+
					"log_level: debug\n"+
					"post_steps:\n"+
					"  go_get_update: false\n"+
					"  git_commit: false\n",
			),
			0o644,
		),
	)

	cfg := Default()
	err := cfg.LoadYAML(configPath)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/src", cfg.RootPath)
	assert.Equal(t, "github.com/acme", cfg.GitLocation)
	assert.Equal(t, zerolog.DebugLevel, cfg.LogLevel)
	assert.False(t, cfg.PostSteps.GoGetUpdate)
	assert.True(t, cfg.PostSteps.GoModTidy)
	assert.True(t, cfg.PostSteps.GitInit)
	assert.False(t, cfg.PostSteps.GitCommit)
}

func TestLoadYAMLExpandsEnvironmentVariables(t *testing.T) {
	t.Setenv("BCLI_TEST_HOME", "/tmp/home")

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(
		t,
		os.WriteFile(
			configPath,
			[]byte("root_path: $BCLI_TEST_HOME/src\ngit_location: github.com/acme\n"),
			0o644,
		),
	)

	cfg := Default()
	err := cfg.LoadYAML(configPath)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/home/src", cfg.RootPath)
	assert.Equal(t, "github.com/acme", cfg.GitLocation)
}
