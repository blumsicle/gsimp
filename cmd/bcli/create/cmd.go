// Package create implements the bcli create subcommand.
package create

import (
	"context"
	"os"
	"path/filepath"

	"github.com/blumsicle/bcli/internal/bcliconfig"
	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/blumsicle/bcli/internal/poststep"
	"github.com/blumsicle/bcli/internal/projectgen"
	"github.com/rs/zerolog"
)

// Command creates a new starter CLI project.
type Command struct {
	RootPath         *string `short:"r" help:"Directory to create the new project under"                                              type:"path" completion:"<directory>" xor:"inplace-root"`
	ProjectDirPrefix *string `short:"p" help:"Prefix to prepend to the generated project directory name"                                                                   xor:"inplace-prefix"`
	GitLocation      *string `short:"g" help:"Git host and owner prefix for the generated module path"`
	NoGoGetUpdate    bool    `          help:"Skip the 'go get -u ./...' post step"`
	NoGoModTidy      bool    `          help:"Skip the 'go mod tidy' post step"`
	NoGitInit        bool    `          help:"Skip the 'git init' post step"`
	NoGitCommit      bool    `          help:"Skip the 'git commit' post step"`
	InPlace          bool    `          help:"Create the project in the current directory, ignoring root_path and project_dir_prefix"                                      xor:"inplace-root,inplace-prefix" name:"inplace"`
	JSON             bool    `          help:"Write project creation metadata as JSON to stdout"                                                                                                             name:"json"`
	Name             string  `          help:"Name of the new CLI project"                                                                                                                                                  arg:"" required:""`
	Description      string  `          help:"Description for the generated CLI"                                                                                                                                            arg:"" required:""`
}

// AfterApply applies command-specific flag overrides to the shared app config.
func (c *Command) AfterApply(cfg *bcliconfig.Config) error {
	cfg.ApplyCreateOverrides(bcliconfig.CreateOverrides{
		RootPath:         c.RootPath,
		ProjectDirPrefix: c.ProjectDirPrefix,
		GitLocation:      c.GitLocation,
		NoGoGetUpdate:    c.NoGoGetUpdate,
		NoGoModTidy:      c.NoGoModTidy,
		NoGitInit:        c.NoGitInit,
		NoGitCommit:      c.NoGitCommit,
	})

	return nil
}

// Run generates the project scaffold and executes the configured post steps.
func (c *Command) Run(log zerolog.Logger, cfg *bcliconfig.Config) error {
	log = cliutil.SubLogger(log, "create")
	cfg.Normalize()
	rootPath := cfg.RootPath
	projectDirPrefix := cfg.ProjectDirPrefix
	if c.InPlace {
		rootPath = ""
		projectDirPrefix = ""
	}
	log.Debug().
		Str("root_path", rootPath).
		Str("git_location", cfg.GitLocation).
		Str("project_dir_prefix", projectDirPrefix).
		Bool("inplace", c.InPlace).
		Bool("go_get_update", cfg.PostSteps.GoGetUpdate).
		Bool("go_mod_tidy", cfg.PostSteps.GoModTidy).
		Bool("git_init", cfg.PostSteps.GitInit).
		Bool("git_commit", cfg.PostSteps.GitCommit).
		Msg("resolved create command configuration")

	gen := projectgen.New(log)
	planner := poststep.NewPlanner(log, &cfg.PostSteps)
	plannedSteps := planner.Planned()
	for _, step := range plannedSteps {
		gen.AddPostStep(step)
	}

	result, err := gen.GenerateResult(context.Background(), projectgen.Config{
		Name:             c.Name,
		Description:      c.Description,
		GitLocation:      cfg.GitLocation,
		ProjectDirPrefix: projectDirPrefix,
		RootPath:         rootPath,
		InPlace:          c.InPlace,
	})
	if err != nil {
		return err
	}

	targetPath := filepath.Clean(result.TargetPath)
	log.Info().
		Str("project", c.Name).
		Str("path", targetPath).
		Msg("created project scaffold")

	if c.JSON {
		createResult, err := NewCreateResult(c.Name, c.Description, c.InPlace, result, plannedSteps)
		if err != nil {
			return err
		}
		if err := WriteCreateJSON(os.Stdout, createResult); err != nil {
			return err
		}
	}

	return nil
}
