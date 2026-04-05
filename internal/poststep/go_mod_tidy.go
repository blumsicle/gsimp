package poststep

import "context"

type GoModTidyPostStep struct{}

func (GoModTidyPostStep) Name() string {
	return "go mod tidy"
}

func (GoModTidyPostStep) Run(ctx context.Context, input PostStepInput) error {
	return run(ctx, input.ProjectPath, "go", "mod", "tidy")
}
