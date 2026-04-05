package poststep

import "context"

type GitInit struct{}

func (GitInit) Name() string {
	return "git init"
}

func (GitInit) Run(ctx context.Context, input Input) error {
	if err := runCommand(ctx, input.ProjectPath, "git", "init"); err != nil {
		return err
	}
	if err := runCommand(ctx, input.ProjectPath, "git", "add", "."); err != nil {
		return err
	}
	return runCommand(
		ctx,
		input.ProjectPath,
		"git",
		"commit",
		"-m",
		"Initial commit",
	)
}
