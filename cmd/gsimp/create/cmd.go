// Package create implements the gsimp create subcommand.
package create

import (
	"context"
	"path/filepath"

	"github.com/blumsicle/gsimp/internal/appconfig"
	"github.com/blumsicle/gsimp/internal/poststep"
	"github.com/blumsicle/gsimp/internal/projectgen"
	"github.com/rs/zerolog"
)

// Command creates a new starter CLI project.
type Command struct {
	RootPath      *string `short:"r" help:"Directory to create the new project under"               type:"path"`
	GitLocation   *string `short:"g" help:"Git host and owner prefix for the generated module path"`
	NoGoGetUpdate bool    `          help:"Skip the 'go get -u ./...' post step"`
	NoGoModTidy   bool    `          help:"Skip the 'go mod tidy' post step"`
	NoGitInit     bool    `          help:"Skip the 'git init' post step"`
	NoGitCommit   bool    `          help:"Skip the 'git commit' post step"`
	Name          string  `          help:"Name of the new CLI project"                                         arg:"" required:""`
	Description   string  `          help:"Description for the generated CLI"                                   arg:"" required:""`
}

// AfterApply applies command-specific flag overrides to the shared app config.
func (c *Command) AfterApply(cfg *appconfig.Config) error {
	if c.RootPath != nil {
		cfg.RootPath = *c.RootPath
	}
	if c.GitLocation != nil {
		cfg.GitLocation = *c.GitLocation
	}
	if c.NoGoGetUpdate {
		cfg.PostSteps.GoGetUpdate = false
	}
	if c.NoGoModTidy {
		cfg.PostSteps.GoModTidy = false
	}
	if c.NoGitInit {
		cfg.PostSteps.GitInit = false
	}
	if c.NoGitCommit {
		cfg.PostSteps.GitCommit = false
	}

	return nil
}

// Run generates the project scaffold and executes the configured post steps.
func (c *Command) Run(log zerolog.Logger, cfg *appconfig.Config) error {
	gen := projectgen.New()
	for _, step := range poststep.Planned(cfg) {
		gen.AddPostStep(step)
	}

	targetPath, err := gen.Generate(context.Background(), projectgen.Config{
		Name:        c.Name,
		Description: c.Description,
		GitLocation: cfg.GitLocation,
		RootPath:    cfg.RootPath,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("project", c.Name).
		Str("path", filepath.Clean(targetPath)).
		Msg("created project")

	return nil
}
