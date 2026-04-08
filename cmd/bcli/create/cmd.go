// Package create implements the bcli create subcommand.
package create

import (
	"context"
	"path/filepath"

	"github.com/blumsicle/bcli/internal/appconfig"
	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/blumsicle/bcli/internal/poststep"
	"github.com/blumsicle/bcli/internal/projectgen"
	"github.com/rs/zerolog"
)

// Command creates a new starter CLI project.
type Command struct {
	RootPath         *string `short:"r" help:"Directory to create the new project under"                 type:"path"`
	ProjectDirPrefix *string `short:"p" help:"Prefix to prepend to the generated project directory name"`
	GitLocation      *string `short:"g" help:"Git host and owner prefix for the generated module path"`
	NoGoGetUpdate    bool    `          help:"Skip the 'go get -u ./...' post step"`
	NoGoModTidy      bool    `          help:"Skip the 'go mod tidy' post step"`
	NoGitInit        bool    `          help:"Skip the 'git init' post step"`
	NoGitCommit      bool    `          help:"Skip the 'git commit' post step"`
	Name             string  `          help:"Name of the new CLI project"                                           arg:"" required:""`
	Description      string  `          help:"Description for the generated CLI"                                     arg:"" required:""`
}

// AfterApply applies command-specific flag overrides to the shared app config.
func (c *Command) AfterApply(cfg *appconfig.Config) error {
	if c.RootPath != nil {
		cfg.RootPath = *c.RootPath
	}
	if c.ProjectDirPrefix != nil {
		cfg.ProjectDirPrefix = *c.ProjectDirPrefix
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
	log = cliutil.SubLogger(log, "create")
	cfg.Normalize()
	log.Debug().
		Str("root_path", cfg.RootPath).
		Str("git_location", cfg.GitLocation).
		Str("project_dir_prefix", cfg.ProjectDirPrefix).
		Bool("go_get_update", cfg.PostSteps.GoGetUpdate).
		Bool("go_mod_tidy", cfg.PostSteps.GoModTidy).
		Bool("git_init", cfg.PostSteps.GitInit).
		Bool("git_commit", cfg.PostSteps.GitCommit).
		Msg("resolved create command configuration")

	gen := projectgen.New(log)
	planner := poststep.NewPlanner(log, &cfg.PostSteps)
	for _, step := range planner.Planned() {
		gen.AddPostStep(step)
	}

	targetPath, err := gen.Generate(context.Background(), projectgen.Config{
		Name:             c.Name,
		Description:      c.Description,
		GitLocation:      cfg.GitLocation,
		ProjectDirPrefix: cfg.ProjectDirPrefix,
		RootPath:         cfg.RootPath,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("project", c.Name).
		Str("path", filepath.Clean(targetPath)).
		Msg("created project scaffold")

	return nil
}
