// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import "context"

// GoGetUpdatePostStep updates module dependencies in the generated project.
type GoGetUpdatePostStep struct{}

// Name returns the human-readable name of the post step.
func (GoGetUpdatePostStep) Name() string {
	return "go get -u ./..."
}

// Run executes `go get -u ./...` in the generated project.
func (GoGetUpdatePostStep) Run(ctx context.Context, input PostStepInput) error {
	return run(ctx, input.ProjectPath, "go", "get", "-u", "./...")
}
