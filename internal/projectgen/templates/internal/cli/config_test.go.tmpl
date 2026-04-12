package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testYAMLConfig struct {
	Name string `yaml:"name"`
}

func TestLoadYAMLMissingFileKeepsConfig(t *testing.T) {
	cfg := testYAMLConfig{Name: "default"}

	err := LoadYAML(filepath.Join(t.TempDir(), "missing.yaml"), &cfg, "test")
	require.NoError(t, err)

	assert.Equal(t, "default", cfg.Name)
}

func TestLoadYAMLOverridesConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte("name: from-file\n"), 0o644))
	cfg := testYAMLConfig{Name: "default"}

	err := LoadYAML(path, &cfg, "test")
	require.NoError(t, err)

	assert.Equal(t, "from-file", cfg.Name)
}

func TestLoadYAMLWrapsReadErrors(t *testing.T) {
	cfg := testYAMLConfig{}

	err := LoadYAML(t.TempDir(), &cfg, "test")

	require.Error(t, err)
	assert.ErrorContains(t, err, "read test config")
}

func TestLoadYAMLWrapsParseErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte("name: [\n"), 0o644))
	cfg := testYAMLConfig{}

	err := LoadYAML(path, &cfg, "test")

	require.Error(t, err)
	assert.ErrorContains(t, err, "parse test config")
}
