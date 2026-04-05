package create

import (
	"context"
	"path/filepath"

	"github.com/blumsicle/gsimp/internal/appconfig"
	"github.com/blumsicle/gsimp/internal/poststep"
	"github.com/blumsicle/gsimp/internal/projectgen"
	"github.com/rs/zerolog"
)

type Command struct {
	RootPath    *string `short:"r" help:"Directory to create the new project under"               type:"path"`
	GitLocation *string `short:"g" help:"Git host and owner prefix for the generated module path"`
	Name        string  `          help:"Name of the new CLI project"                                         arg:"" required:""`
	Description string  `          help:"Description for the generated CLI"                                   arg:"" required:""`
}

func (c *Command) AfterApply(cfg *appconfig.Config) error {
	if c.RootPath != nil {
		cfg.RootPath = *c.RootPath
	}
	if c.GitLocation != nil {
		cfg.GitLocation = *c.GitLocation
	}

	return nil
}

func (c *Command) Run(log zerolog.Logger, cfg *appconfig.Config) error {
	gen := projectgen.New()
	for _, step := range poststep.DefaultPostSteps() {
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
