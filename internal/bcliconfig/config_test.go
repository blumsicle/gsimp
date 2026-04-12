package bcliconfig

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	assert.Equal(t, ".", cfg.RootPath)
	assert.Equal(t, "", cfg.ProjectDirPrefix)
	assert.Equal(t, "", cfg.GitLocation)
	assert.Equal(t, zerolog.InfoLevel, cfg.LogLevel)
	assert.True(t, cfg.PostSteps.GoGetUpdate)
	assert.True(t, cfg.PostSteps.GoModTidy)
	assert.True(t, cfg.PostSteps.GitInit)
	assert.True(t, cfg.PostSteps.GitCommit)
}
