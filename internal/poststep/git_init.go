// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import (
	"context"

	"github.com/rs/zerolog"
)

// GitInitPostStep initializes a Git repository in the generated project.
type GitInitPostStep struct {
	log zerolog.Logger
}

// Name returns the human-readable name of the post step.
func (GitInitPostStep) Name() string {
	return "git init"
}

// Run initializes Git in the generated project.
func (s GitInitPostStep) Run(ctx context.Context, input PostStepInput) error {
	s.log.Info().Str("project_path", input.ProjectPath).Msg("initializing git repository")
	return run(ctx, s.log, input.ProjectPath, "git", "init")
}
