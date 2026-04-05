package poststep

import "context"

type GoGetUpdatePostStep struct{}

func (GoGetUpdatePostStep) Name() string {
	return "go get -u ./..."
}

func (GoGetUpdatePostStep) Run(ctx context.Context, input PostStepInput) error {
	return run(ctx, input.ProjectPath, "go", "get", "-u", "./...")
}
