package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/blumsicle/bcli/internal/appconfig"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStdout(t *testing.T) (*os.File, func() []byte) {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = writer

	return writer, func() []byte {
		require.NoError(t, writer.Close())
		os.Stdout = originalStdout

		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.NoError(t, reader.Close())

		return data
	}
}

func TestRunWritesResolvedConfigToStdout(t *testing.T) {
	command := Command{}
	cfg := appconfig.Default()
	cfg.RootPath = "/tmp/src"
	cfg.GitLocation = "github.com/acme"
	cfg.PostSteps.GitCommit = false

	_, restoreStdout := captureStdout(t)

	err := command.Run(zerolog.Nop(), cfg)
	require.NoError(t, err)

	var got appconfig.Config
	require.NoError(t, yaml.Unmarshal(restoreStdout(), &got))
	assert.Equal(t, *cfg, got)
}

func TestRunWritesResolvedConfigToFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "resolved.yaml")
	command := Command{Output: outputPath}
	cfg := appconfig.Default()
	cfg.RootPath = "/tmp/src"
	cfg.GitLocation = "github.com/acme"
	cfg.PostSteps.GoGetUpdate = false

	err := command.Run(zerolog.Nop(), cfg)
	require.NoError(t, err)

	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	var got appconfig.Config
	require.NoError(t, yaml.Unmarshal(data, &got))
	assert.Equal(t, *cfg, got)
}
