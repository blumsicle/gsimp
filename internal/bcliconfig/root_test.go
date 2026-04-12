package bcliconfig

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestApplyRootOverrides(t *testing.T) {
	cfg := Default()
	cfg.LogLevel = zerolog.DebugLevel

	cfg.ApplyRootOverrides(RootOverrides{})

	assert.Equal(t, zerolog.DebugLevel, cfg.LogLevel)

	logLevel := zerolog.WarnLevel
	cfg.ApplyRootOverrides(RootOverrides{
		LogLevel: &logLevel,
	})

	assert.Equal(t, zerolog.WarnLevel, cfg.LogLevel)
}
