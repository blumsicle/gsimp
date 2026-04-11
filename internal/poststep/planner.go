// Package poststep defines post-generation steps run after a scaffold is written.
package poststep

import (
	"github.com/blumsicle/bcli/internal/appconfig"
	cliutil "github.com/blumsicle/bcli/internal/cli"
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

type definitionSpec struct {
	id       StepID
	requires []StepID
}

var postStepDefinitions = []definitionSpec{
	{id: StepGoGetUpdate},
	{id: StepGoModTidy},
	{id: StepGitInit},
	{id: StepGitCommit, requires: []StepID{StepGitInit}},
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

// Planned returns the enabled post steps in execution order after applying dependencies.
func (p *Planner) Planned() []PostStep {
	definitions := postStepDefinitions
	enabled := make(map[StepID]bool, len(definitions))

	for _, definition := range definitions {
		enabled[definition.id] = definition.enabled(p.cfg)
		p.log.Debug().
			Str("step_id", string(definition.id)).
			Bool("enabled", enabled[definition.id]).
			Msg("evaluated post-step config")
	}

	steps := make([]PostStep, 0, len(definitions))
	for _, definition := range definitions {
		if !enabled[definition.id] {
			p.log.Debug().
				Str("step_id", string(definition.id)).
				Msg("skipping disabled post step")
			continue
		}

		satisfied := true
		for _, dependency := range definition.requires {
			if !enabled[dependency] {
				satisfied = false
				p.log.Debug().
					Str("step_id", string(definition.id)).
					Str("dependency", string(dependency)).
					Msg("skipping post step because dependency is disabled")
				break
			}
		}
		if !satisfied {
			continue
		}

		step := definition.postStep(p.log)
		p.log.Info().
			Str("step_id", string(definition.id)).
			Str("name", step.Name()).
			Msg("selected post step")
		steps = append(steps, step)
	}

	return steps
}

func (d definitionSpec) enabled(cfg *appconfig.PostStepsConfig) bool {
	switch d.id {
	case StepGoGetUpdate:
		return cfg.GoGetUpdate
	case StepGoModTidy:
		return cfg.GoModTidy
	case StepGitInit:
		return cfg.GitInit
	case StepGitCommit:
		return cfg.GitCommit
	default:
		return false
	}
}

func (d definitionSpec) postStep(log zerolog.Logger) PostStep {
	switch d.id {
	case StepGoGetUpdate:
		return GoGetUpdatePostStep{log: log}
	case StepGoModTidy:
		return GoModTidyPostStep{log: log}
	case StepGitInit:
		return GitInitPostStep{log: log}
	case StepGitCommit:
		return GitCommitPostStep{log: log}
	default:
		return nil
	}
}
