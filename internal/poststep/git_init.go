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
	return gitInitSpec.name
}

// Run initializes Git in the generated project.
func (s GitInitPostStep) Run(ctx context.Context, input PostStepInput) error {
	return gitInitSpec.run(ctx, s.log, input)
}

var gitInitSpec = commandPostStepSpec{
	name:    "git init",
	message: "initializing git repository",
	command: "git",
	args:    []string{"init"},
}
