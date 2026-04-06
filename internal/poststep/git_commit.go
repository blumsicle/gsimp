// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import (
	"context"

	"github.com/rs/zerolog"
)

// GitCommitPostStep stages the generated scaffold and creates the initial commit.
type GitCommitPostStep struct {
	log zerolog.Logger
}

// Name returns the human-readable name of the post step.
func (GitCommitPostStep) Name() string {
	return "git commit"
}

// Run stages all generated files and creates the initial commit.
func (s GitCommitPostStep) Run(ctx context.Context, input PostStepInput) error {
	s.log.Info().Str("project_path", input.ProjectPath).Msg("creating initial git commit")
	if err := run(ctx, s.log, input.ProjectPath, "git", "add", "."); err != nil {
		return err
	}

	return run(
		ctx,
		s.log,
		input.ProjectPath,
		"git",
		"commit",
		"-m",
		"Initial commit",
	)
}
