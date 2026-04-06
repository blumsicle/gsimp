// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import "context"

// GitInitPostStep initializes a Git repository in the generated project.
type GitInitPostStep struct{}

// Name returns the human-readable name of the post step.
func (GitInitPostStep) Name() string {
	return "git init"
}

// Run initializes Git in the generated project.
func (GitInitPostStep) Run(ctx context.Context, input PostStepInput) error {
	return run(ctx, input.ProjectPath, "git", "init")
}
