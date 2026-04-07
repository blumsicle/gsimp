// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import (
	"github.com/blumsicle/bcli/internal/appconfig"
	"github.com/rs/zerolog"
)

// DefaultPostSteps returns the default ordered post steps for a generated project.
func DefaultPostSteps() []PostStep {
	cfg := appconfig.Default()
	return NewPlanner(zerolog.Nop(), &cfg.PostSteps).Planned()
}
