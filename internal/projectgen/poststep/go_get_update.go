package poststep

import "context"

type GoGetUpdate struct{}

func (GoGetUpdate) Name() string {
	return "go get -u ./..."
}

func (GoGetUpdate) Run(ctx context.Context, input Input) error {
	return runCommand(ctx, input.ProjectPath, "go", "get", "-u", "./...")
}
