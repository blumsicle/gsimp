package create

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/blumsicle/bcli/internal/poststep"
	"github.com/blumsicle/bcli/internal/projectgen"
)

// CreateResult describes a project created by the create command.
type CreateResult struct {
	Project     string           `json:"project"`
	Description string           `json:"description"`
	ModulePath  string           `json:"module_path"`
	TargetPath  string           `json:"target_path"`
	InPlace     bool             `json:"inplace"`
	PostSteps   []PostStepResult `json:"post_steps"`
}

// PostStepResult describes whether a create command post step ran.
type PostStepResult struct {
	Name string `json:"name"`
	Ran  bool   `json:"ran"`
}

// NewCreateResult returns the JSON-serializable result for a generated project.
func NewCreateResult(
	name string,
	description string,
	inPlace bool,
	result projectgen.Result,
	plannedSteps []poststep.PostStep,
) (CreateResult, error) {
	absoluteTargetPath, err := filepath.Abs(result.TargetPath)
	if err != nil {
		return CreateResult{}, fmt.Errorf("resolve target path: %w", err)
	}

	return CreateResult{
		Project:     name,
		Description: description,
		ModulePath:  result.ModulePath,
		TargetPath:  absoluteTargetPath,
		InPlace:     inPlace,
		PostSteps:   PostStepResults(plannedSteps),
	}, nil
}

// WriteCreateJSON writes a create command result as JSON.
func WriteCreateJSON(w io.Writer, result CreateResult) error {
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return fmt.Errorf("write create json: %w", err)
	}

	return nil
}

// PostStepResults returns the standard create command post-step result set.
func PostStepResults(plannedSteps []poststep.PostStep) []PostStepResult {
	ran := make(map[string]bool, len(plannedSteps))
	for _, step := range plannedSteps {
		ran[step.Name()] = true
	}

	return []PostStepResult{
		{Name: "go get -u ./...", Ran: ran["go get -u ./..."]},
		{Name: "go mod tidy", Ran: ran["go mod tidy"]},
		{Name: "git init", Ran: ran["git init"]},
		{Name: "git commit", Ran: ran["git commit"]},
	}
}
