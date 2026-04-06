// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import (
	"context"

	"github.com/rs/zerolog"
)

// GoGetUpdatePostStep updates module dependencies in the generated project.
type GoGetUpdatePostStep struct {
	log zerolog.Logger
}

// Name returns the human-readable name of the post step.
func (GoGetUpdatePostStep) Name() string {
	return "go get -u ./..."
}

// Run executes `go get -u ./...` in the generated project.
func (s GoGetUpdatePostStep) Run(ctx context.Context, input PostStepInput) error {
	s.log.Info().Str("project_path", input.ProjectPath).Msg("updating module dependencies")
	return run(ctx, s.log, input.ProjectPath, "go", "get", "-u", "./...")
}
