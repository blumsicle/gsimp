package poststep

import "context"

type PostStepInput struct {
	ProjectPath string
	Name        string
	ModulePath  string
}

type PostStep interface {
	Name() string
	Run(ctx context.Context, input PostStepInput) error
}
