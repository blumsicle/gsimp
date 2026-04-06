// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import (
	"context"

	"github.com/rs/zerolog"
)

// GoModTidyPostStep tidies the generated project's module metadata.
type GoModTidyPostStep struct {
	log zerolog.Logger
}

// Name returns the human-readable name of the post step.
func (GoModTidyPostStep) Name() string {
	return "go mod tidy"
}

// Run executes `go mod tidy` in the generated project.
func (s GoModTidyPostStep) Run(ctx context.Context, input PostStepInput) error {
	s.log.Info().Str("project_path", input.ProjectPath).Msg("tidying module metadata")
	return run(ctx, s.log, input.ProjectPath, "go", "mod", "tidy")
}
