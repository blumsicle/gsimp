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
	return goGetUpdateSpec.name
}

// Run executes `go get -u ./...` in the generated project.
func (s GoGetUpdatePostStep) Run(ctx context.Context, input PostStepInput) error {
	return goGetUpdateSpec.run(ctx, s.log, input)
}

var goGetUpdateSpec = commandPostStepSpec{
	name:    "go get -u ./...",
	message: "updating module dependencies",
	command: "go",
	args:    []string{"get", "-u", "./..."},
}
