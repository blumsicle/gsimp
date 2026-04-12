package bcliconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyCreateOverrides(t *testing.T) {
	rootPath := "/tmp/src"
	projectDirPrefix := "generated-"
	gitLocation := "github.com/acme"
	cfg := Default()

	cfg.ApplyCreateOverrides(CreateOverrides{
		RootPath:         &rootPath,
		ProjectDirPrefix: &projectDirPrefix,
		GitLocation:      &gitLocation,
		NoGoGetUpdate:    true,
		NoGoModTidy:      true,
		NoGitInit:        true,
		NoGitCommit:      true,
	})

	assert.Equal(t, rootPath, cfg.RootPath)
	assert.Equal(t, projectDirPrefix, cfg.ProjectDirPrefix)
	assert.Equal(t, gitLocation, cfg.GitLocation)
	assert.False(t, cfg.PostSteps.GoGetUpdate)
	assert.False(t, cfg.PostSteps.GoModTidy)
	assert.False(t, cfg.PostSteps.GitInit)
	assert.False(t, cfg.PostSteps.GitCommit)
}

func TestApplyCreateOverridesPreservesExistingValuesWhenOmitted(t *testing.T) {
	cfg := Default()
	cfg.RootPath = "/existing-root"
	cfg.ProjectDirPrefix = "existing-"
	cfg.GitLocation = "github.com/existing"
	cfg.PostSteps.GoGetUpdate = false
	cfg.PostSteps.GoModTidy = false
	cfg.PostSteps.GitInit = false
	cfg.PostSteps.GitCommit = false

	cfg.ApplyCreateOverrides(CreateOverrides{})

	assert.Equal(t, "/existing-root", cfg.RootPath)
	assert.Equal(t, "existing-", cfg.ProjectDirPrefix)
	assert.Equal(t, "github.com/existing", cfg.GitLocation)
	assert.False(t, cfg.PostSteps.GoGetUpdate)
	assert.False(t, cfg.PostSteps.GoModTidy)
	assert.False(t, cfg.PostSteps.GitInit)
	assert.False(t, cfg.PostSteps.GitCommit)
}
