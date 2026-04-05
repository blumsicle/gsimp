// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import "context"

// PostStepInput describes the generated project passed to a post step.
type PostStepInput struct {
	ProjectPath string
	Name        string
	ModulePath  string
}

// PostStep represents a single post-generation action.
type PostStep interface {
	Name() string
	Run(ctx context.Context, input PostStepInput) error
}
