// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import (
	"github.com/blumsicle/gsimp/internal/appconfig"
	cliutil "github.com/blumsicle/gsimp/internal/cli"
	"github.com/rs/zerolog"
)

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
	Enabled  func(cfg *appconfig.PostStepsConfig) bool
}

// Planner builds the ordered post-step plan for a generator run.
type Planner struct {
	log zerolog.Logger
	cfg *appconfig.PostStepsConfig
}

// NewPlanner constructs a Planner with subsystem-specific loggers.
func NewPlanner(log zerolog.Logger, cfg *appconfig.PostStepsConfig) *Planner {
	return &Planner{
		log: cliutil.SubLogger(log, "poststep"),
		cfg: cfg,
	}
}

// Definitions returns the ordered post-step definitions used by the generator.
func (p *Planner) Definitions() []Definition {
	return []Definition{
		{
			ID:       StepGoGetUpdate,
			PostStep: GoGetUpdatePostStep{log: p.log},
			Enabled: func(cfg *appconfig.PostStepsConfig) bool {
				return cfg.GoGetUpdate
			},
		},
		{
			ID:       StepGoModTidy,
			PostStep: GoModTidyPostStep{log: p.log},
			Enabled: func(cfg *appconfig.PostStepsConfig) bool {
				return cfg.GoModTidy
			},
		},
		{
			ID:       StepGitInit,
			PostStep: GitInitPostStep{log: p.log},
			Enabled: func(cfg *appconfig.PostStepsConfig) bool {
				return cfg.GitInit
			},
		},
		{
			ID:       StepGitCommit,
			PostStep: GitCommitPostStep{log: p.log},
			Requires: []StepID{StepGitInit},
			Enabled: func(cfg *appconfig.PostStepsConfig) bool {
				return cfg.GitCommit
			},
		},
	}
}

// Planned returns the enabled post steps in execution order after applying dependencies.
func (p *Planner) Planned() []PostStep {
	definitions := p.Definitions()
	enabled := make(map[StepID]bool, len(definitions))

	for _, definition := range definitions {
		enabled[definition.ID] = definition.Enabled(p.cfg)
		p.log.Debug().
			Str("step_id", string(definition.ID)).
			Bool("enabled", enabled[definition.ID]).
			Msg("evaluated post-step config")
	}

	steps := make([]PostStep, 0, len(definitions))
	for _, definition := range definitions {
		if !enabled[definition.ID] {
			p.log.Debug().
				Str("step_id", string(definition.ID)).
				Msg("skipping disabled post step")
			continue
		}

		satisfied := true
		for _, dependency := range definition.Requires {
			if !enabled[dependency] {
				satisfied = false
				p.log.Debug().
					Str("step_id", string(definition.ID)).
					Str("dependency", string(dependency)).
					Msg("skipping post step because dependency is disabled")
				break
			}
		}
		if !satisfied {
			continue
		}

		p.log.Info().
			Str("step_id", string(definition.ID)).
			Str("name", definition.PostStep.Name()).
			Msg("selected post step")
		steps = append(steps, definition.PostStep)
	}

	return steps
}
