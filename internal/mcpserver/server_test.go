package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/blumsicle/bcli/cmd/bcli/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRunner struct {
	request CommandRequest
	result  CommandResult
	err     error
}

func (r *fakeRunner) Run(_ context.Context, request CommandRequest) (CommandResult, error) {
	r.request = request
	return r.result, r.err
}

func TestCreateGoCLIProjectMapsInputToBCLIArgs(t *testing.T) {
	stdout, err := json.Marshal(createProjectOutput{
		CreateResult: create.CreateResult{
			Project:     "cooltool",
			Description: "CLI tool that does cool stuff",
			ModulePath:  "github.com/acme/cooltool",
			TargetPath:  "/tmp/cooltool",
			InPlace:     true,
			PostSteps: []create.PostStepResult{
				{Name: "go mod tidy", Ran: false},
			},
		},
	})
	require.NoError(t, err)

	runner := &fakeRunner{
		result: CommandResult{
			Stdout: string(stdout),
			Stderr: "created\n",
		},
	}
	server := New(Config{BCLICommand: "bcli", Timeout: time.Minute}, runner)

	got, err := server.CreateGoCLIProject(context.Background(), createProjectInput{
		Name:             "cooltool",
		Description:      "CLI tool that does cool stuff",
		WorkingDirectory: "/work/cooltool",
		RootPath:         "/ignored/when/inplace",
		ProjectDirPrefix: "local-",
		GitLocation:      "github.com/acme",
		BCLIConfigFile:   "/tmp/bcli.yaml",
		InPlace:          true,
		SkipGoGetUpdate:  true,
		SkipGoModTidy:    true,
		SkipGitInit:      true,
		SkipGitCommit:    true,
	})
	require.NoError(t, err)

	assert.Equal(t, "bcli", runner.request.Command)
	assert.Equal(t, "/work/cooltool", runner.request.WorkingDir)
	assert.Equal(
		t,
		[]string{
			"--config-file", "/tmp/bcli.yaml",
			"create", "--json",
			"--git-location", "github.com/acme",
			"--inplace",
			"--no-go-get-update",
			"--no-go-mod-tidy",
			"--no-git-init",
			"--no-git-commit",
			"cooltool",
			"CLI tool that does cool stuff",
		},
		runner.request.Args,
	)
	assert.Equal(t, "cooltool", got.Project)
	assert.Equal(t, "/tmp/cooltool", got.TargetPath)
	assert.Equal(
		t,
		[]string{
			"bcli",
			"--config-file",
			"/tmp/bcli.yaml",
			"create",
			"--json",
			"--git-location",
			"github.com/acme",
			"--inplace",
			"--no-go-get-update",
			"--no-go-mod-tidy",
			"--no-git-init",
			"--no-git-commit",
			"cooltool",
			"CLI tool that does cool stuff",
		},
		got.Command,
	)
	assert.Equal(t, "/work/cooltool", got.WorkingDir)
	assert.Equal(t, "created\n", got.Stderr)
}

func TestCreateGoCLIProjectRequiresName(t *testing.T) {
	server := New(Config{}, &fakeRunner{})

	_, err := server.CreateGoCLIProject(context.Background(), createProjectInput{
		Description: "CLI tool that does cool stuff",
	})

	require.Error(t, err)
	assert.EqualError(t, err, "name is required")
}

func TestCreateGoCLIProjectRequiresDescription(t *testing.T) {
	server := New(Config{}, &fakeRunner{})

	_, err := server.CreateGoCLIProject(context.Background(), createProjectInput{
		Name: "cooltool",
	})

	require.Error(t, err)
	assert.EqualError(t, err, "description is required")
}

func TestCreateGoCLIProjectReportsMissingBCLI(t *testing.T) {
	server := New(Config{}, &fakeRunner{err: errors.New("executable file not found")})

	_, err := server.CreateGoCLIProject(context.Background(), createProjectInput{
		Name:        "cooltool",
		Description: "CLI tool that does cool stuff",
	})

	require.Error(t, err)
	assert.EqualError(t, err, "run bcli: executable file not found")
}

func TestCreateGoCLIProjectReportsNonZeroExit(t *testing.T) {
	server := New(Config{}, &fakeRunner{
		result: CommandResult{
			Stderr:   "target directory is not empty\n",
			ExitCode: 1,
		},
	})

	_, err := server.CreateGoCLIProject(context.Background(), createProjectInput{
		Name:        "cooltool",
		Description: "CLI tool that does cool stuff",
	})

	require.Error(t, err)
	assert.EqualError(t, err, "bcli exited with code 1: target directory is not empty")
}

func TestCreateGoCLIProjectReportsInvalidJSON(t *testing.T) {
	server := New(Config{}, &fakeRunner{
		result: CommandResult{Stdout: "not json"},
	})

	_, err := server.CreateGoCLIProject(context.Background(), createProjectInput{
		Name:        "cooltool",
		Description: "CLI tool that does cool stuff",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "parse bcli json output")
}
