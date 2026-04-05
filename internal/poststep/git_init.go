package poststep

import "context"

type GitInitPostStep struct{}

func (GitInitPostStep) Name() string {
	return "git init"
}

func (GitInitPostStep) Run(ctx context.Context, input PostStepInput) error {
	if err := run(ctx, input.ProjectPath, "git", "init"); err != nil {
		return err
	}
	if err := run(ctx, input.ProjectPath, "git", "add", "."); err != nil {
		return err
	}
	return run(
		ctx,
		input.ProjectPath,
		"git",
		"commit",
		"-m",
		"Initial commit",
	)
}
