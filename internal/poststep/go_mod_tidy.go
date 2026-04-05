// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import "context"

// GoModTidyPostStep tidies the generated project's module metadata.
type GoModTidyPostStep struct{}

// Name returns the human-readable name of the post step.
func (GoModTidyPostStep) Name() string {
	return "go mod tidy"
}

// Run executes `go mod tidy` in the generated project.
func (GoModTidyPostStep) Run(ctx context.Context, input PostStepInput) error {
	return run(ctx, input.ProjectPath, "go", "mod", "tidy")
}
