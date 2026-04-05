package poststep

import "context"

type GoModTidy struct{}

func (GoModTidy) Name() string {
	return "go mod tidy"
}

func (GoModTidy) Run(ctx context.Context, input Input) error {
	return runCommand(ctx, input.ProjectPath, "go", "mod", "tidy")
}
