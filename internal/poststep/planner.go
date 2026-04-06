// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import "github.com/blumsicle/gsimp/internal/appconfig"

// StepID identifies a configurable post-generation step.
type StepID string

const (
	// StepGoGetUpdate updates generated project dependencies.
	StepGoGetUpdate StepID = "go_get_update"
	// StepGoModTidy tidies the generated project's module metadata.
	StepGoModTidy StepID = "go_mod_tidy"
	// StepGitInit initializes Git in the generated project.
	StepGitInit StepID = "git_init"
	// StepGitCommit stages and commits the generated project.
	StepGitCommit StepID = "git_commit"
)

// Definition describes one post-generation step and its planning rules.
type Definition struct {
	ID       StepID
	PostStep PostStep
	Requires []StepID
	Enabled  func(cfg *appconfig.Config) bool
}

// Definitions returns the ordered post-step definitions used by the generator.
func Definitions() []Definition {
	return []Definition{
		{
			ID:       StepGoGetUpdate,
			PostStep: GoGetUpdatePostStep{},
			Enabled: func(cfg *appconfig.Config) bool {
				return cfg.PostSteps.GoGetUpdate
			},
		},
		{
			ID:       StepGoModTidy,
			PostStep: GoModTidyPostStep{},
			Enabled: func(cfg *appconfig.Config) bool {
				return cfg.PostSteps.GoModTidy
			},
		},
		{
			ID:       StepGitInit,
			PostStep: GitInitPostStep{},
			Enabled: func(cfg *appconfig.Config) bool {
				return cfg.PostSteps.GitInit
			},
		},
		{
			ID:       StepGitCommit,
			PostStep: GitCommitPostStep{},
			Requires: []StepID{StepGitInit},
			Enabled: func(cfg *appconfig.Config) bool {
				return cfg.PostSteps.GitCommit
			},
		},
	}
}

// Planned returns the enabled post steps in execution order after applying dependencies.
func Planned(cfg *appconfig.Config) []PostStep {
	definitions := Definitions()
	enabled := make(map[StepID]bool, len(definitions))

	for _, definition := range definitions {
		enabled[definition.ID] = definition.Enabled(cfg)
	}

	steps := make([]PostStep, 0, len(definitions))
	for _, definition := range definitions {
		if !enabled[definition.ID] {
			continue
		}

		satisfied := true
		for _, dependency := range definition.Requires {
			if !enabled[dependency] {
				satisfied = false
				break
			}
		}
		if !satisfied {
			continue
		}

		steps = append(steps, definition.PostStep)
	}

	return steps
}
