// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import (
	"context"

	"github.com/rs/zerolog"
)

type commandPostStepSpec struct {
	name    string
	message string
	command string
	args    []string
}

func (s commandPostStepSpec) run(
	ctx context.Context,
	log zerolog.Logger,
	input PostStepInput,
) error {
	log.Info().Str("project_path", input.ProjectPath).Msg(s.message)
	return run(ctx, log, input.ProjectPath, s.command, s.args...)
}
