package poststep

import "context"

type Input struct {
	ProjectPath string
	Name        string
	ModulePath  string
}

type Step interface {
	Name() string
	Run(ctx context.Context, input Input) error
}
