// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import "context"

// GitCommitPostStep stages the generated scaffold and creates the initial commit.
type GitCommitPostStep struct{}

// Name returns the human-readable name of the post step.
func (GitCommitPostStep) Name() string {
	return "git commit"
}

// Run stages all generated files and creates the initial commit.
func (GitCommitPostStep) Run(ctx context.Context, input PostStepInput) error {
	if err := run(ctx, input.ProjectPath, "git", "add", "."); err != nil {
		return err
	}

	return run(
		ctx,
		input.ProjectPath,
		"git",
		"commit",
		"-m",
		"Initial commit",
	)
}
