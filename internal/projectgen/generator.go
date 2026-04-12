// Package projectgen renders starter project scaffolds from embedded templates.
package projectgen

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	cliutil "github.com/blumsicle/bcli/internal/cli"
	"github.com/blumsicle/bcli/internal/poststep"
	"github.com/rs/zerolog"
)

// Generator renders embedded project templates and runs post-generation steps.
type Generator struct {
	templateFS fs.FS
	postSteps  []poststep.PostStep
	log        zerolog.Logger
}

// Result describes a completed project generation run.
type Result struct {
	TargetPath string
	ModulePath string
}

// New constructs a generator with the embedded template filesystem.
func New(log zerolog.Logger) *Generator {
	return &Generator{
		templateFS: templateFS,
		postSteps:  []poststep.PostStep{},
		log:        cliutil.SubLogger(log, "projectgen"),
	}
}

// AddPostStep appends a post-generation step to be run after rendering completes.
func (g *Generator) AddPostStep(step poststep.PostStep) {
	g.log.Debug().Str("step", step.Name()).Msg("registered post step")
	g.postSteps = append(g.postSteps, step)
}

// Generate renders the scaffold into the target directory and runs post steps.
func (g *Generator) Generate(ctx context.Context, cfg Config) (string, error) {
	result, err := g.GenerateResult(ctx, cfg)
	if err != nil {
		return "", err
	}

	return result.TargetPath, nil
}

// GenerateResult renders the scaffold and returns detailed generation metadata.
func (g *Generator) GenerateResult(ctx context.Context, cfg Config) (Result, error) {
	if err := validateConfig(cfg); err != nil {
		return Result{}, err
	}

	plan := newGenerationPlan(cfg)
	g.log.Info().
		Str("name", cfg.Name).
		Str("project_dir_prefix", cfg.ProjectDirPrefix).
		Str("root_path", filepath.Clean(plan.rootPath)).
		Str("target_path", filepath.Clean(plan.targetPath)).
		Msg("starting project generation")

	g.log.Debug().
		Str("name", cfg.Name).
		Str("module_path", plan.modulePath).
		Str("go_version", plan.templateData.GoVersion).
		Msg("resolved module path")

	if err := ensureTargetDir(plan.targetPath, cfg.InPlace); err != nil {
		return Result{}, err
	}
	if err := g.renderTemplates(plan.targetPath, plan.templateData); err != nil {
		return Result{}, err
	}
	if err := g.runPostSteps(ctx, plan.postStepInput); err != nil {
		return Result{}, err
	}

	g.log.Info().
		Str("name", cfg.Name).
		Str("target_path", filepath.Clean(plan.targetPath)).
		Msg("finished project generation")

	return Result{
		TargetPath: plan.targetPath,
		ModulePath: plan.modulePath,
	}, nil
}

func (g *Generator) runPostSteps(ctx context.Context, input poststep.PostStepInput) error {
	for _, step := range g.postSteps {
		g.log.Info().Str("step", step.Name()).Msg("running post step")
		if err := step.Run(ctx, input); err != nil {
			return fmt.Errorf("run post step %q: %w", step.Name(), err)
		}
	}

	return nil
}
