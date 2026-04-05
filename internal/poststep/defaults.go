// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

// DefaultPostSteps returns the default ordered post steps for a generated project.
func DefaultPostSteps() []PostStep {
	return []PostStep{
		GoGetUpdatePostStep{},
		GoModTidyPostStep{},
		GitInitPostStep{},
	}
}
