package bcliconfig

// CreateOverrides contains create command flag values that can override config.
type CreateOverrides struct {
	RootPath         *string
	ProjectDirPrefix *string
	GitLocation      *string
	NoGoGetUpdate    bool
	NoGoModTidy      bool
	NoGitInit        bool
	NoGitCommit      bool
}

// ApplyCreateOverrides applies create command flag overrides to the config.
func (c *Config) ApplyCreateOverrides(overrides CreateOverrides) {
	if overrides.RootPath != nil {
		c.RootPath = *overrides.RootPath
	}
	if overrides.ProjectDirPrefix != nil {
		c.ProjectDirPrefix = *overrides.ProjectDirPrefix
	}
	if overrides.GitLocation != nil {
		c.GitLocation = *overrides.GitLocation
	}
	if overrides.NoGoGetUpdate {
		c.PostSteps.GoGetUpdate = false
	}
	if overrides.NoGoModTidy {
		c.PostSteps.GoModTidy = false
	}
	if overrides.NoGitInit {
		c.PostSteps.GitInit = false
	}
	if overrides.NoGitCommit {
		c.PostSteps.GitCommit = false
	}
}
