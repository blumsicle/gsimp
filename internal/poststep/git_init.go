// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import "context"

// GitInitPostStep initializes a Git repository and creates the initial commit.
type GitInitPostStep struct{}

// Name returns the human-readable name of the post step.
func (GitInitPostStep) Name() string {
	return "git init"
}

// Run initializes Git, stages the scaffold, and creates the initial commit.
func (GitInitPostStep) Run(ctx context.Context, input PostStepInput) error {
	if err := run(ctx, input.ProjectPath, "git", "init"); err != nil {
		return err
	}
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
