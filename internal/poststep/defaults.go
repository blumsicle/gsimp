// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import "github.com/blumsicle/gsimp/internal/appconfig"

// DefaultPostSteps returns the default ordered post steps for a generated project.
func DefaultPostSteps() []PostStep {
	return Planned(appconfig.Default())
}
