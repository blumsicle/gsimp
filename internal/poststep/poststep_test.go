package poststep

import (
	"context"
	"errors"
	"testing"

	"github.com/blumsicle/bcli/internal/appconfig"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPostSteps(t *testing.T) {
	steps := DefaultPostSteps()

	require.Len(t, steps, 4)
	assert.IsType(t, GoGetUpdatePostStep{}, steps[0])
	assert.IsType(t, GoModTidyPostStep{}, steps[1])
	assert.IsType(t, GitInitPostStep{}, steps[2])
	assert.IsType(t, GitCommitPostStep{}, steps[3])
}

func TestGoGetUpdatePostStepRunsExpectedCommand(t *testing.T) {
	var (
		dir  string
		name string
		args []string
	)
	previousRun := run
	run = func(_ context.Context, _ zerolog.Logger, gotDir string, gotName string, gotArgs ...string) error {
		dir = gotDir
		name = gotName
		args = gotArgs
		return nil
	}
	t.Cleanup(func() {
		run = previousRun
	})

	err := GoGetUpdatePostStep{log: zerolog.Nop()}.Run(
		context.Background(),
		PostStepInput{ProjectPath: "/tmp/project"},
	)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/project", dir)
	assert.Equal(t, "go", name)
	assert.Equal(t, []string{"get", "-u", "./..."}, args)
}

func TestGoModTidyPostStepRunsExpectedCommand(t *testing.T) {
	var (
		dir  string
		name string
		args []string
	)
	previousRun := run
	run = func(_ context.Context, _ zerolog.Logger, gotDir string, gotName string, gotArgs ...string) error {
		dir = gotDir
		name = gotName
		args = gotArgs
		return nil
	}
	t.Cleanup(func() {
		run = previousRun
	})

	err := GoModTidyPostStep{log: zerolog.Nop()}.Run(
		context.Background(),
		PostStepInput{ProjectPath: "/tmp/project"},
	)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/project", dir)
	assert.Equal(t, "go", name)
	assert.Equal(t, []string{"mod", "tidy"}, args)
}

func TestPlannedDisablesGitCommitWhenGitInitIsDisabled(t *testing.T) {
	cfg := appconfig.Default()
	cfg.PostSteps.GitInit = false
	cfg.PostSteps.GitCommit = true

	steps := NewPlanner(zerolog.Nop(), &cfg.PostSteps).Planned()

	require.Len(t, steps, 2)
	assert.Equal(
		t,
		[]string{"go get -u ./...", "go mod tidy"},
		[]string{steps[0].Name(), steps[1].Name()},
	)
}

func TestGitInitPostStepRunsExpectedCommand(t *testing.T) {
	var calls [][]string
	previousRun := run
	run = func(_ context.Context, _ zerolog.Logger, gotDir string, gotName string, gotArgs ...string) error {
		call := []string{gotDir, gotName}
		call = append(call, gotArgs...)
		calls = append(calls, call)
		return nil
	}
	t.Cleanup(func() {
		run = previousRun
	})

	err := GitInitPostStep{log: zerolog.Nop()}.Run(
		context.Background(),
		PostStepInput{ProjectPath: "/tmp/project"},
	)
	require.NoError(t, err)

	assert.Equal(
		t,
		[][]string{
			{"/tmp/project", "git", "init"},
		},
		calls,
	)
}

func TestGitCommitPostStepRunsExpectedCommandsInOrder(t *testing.T) {
	var calls [][]string
	previousRun := run
	run = func(_ context.Context, _ zerolog.Logger, gotDir string, gotName string, gotArgs ...string) error {
		call := []string{gotDir, gotName}
		call = append(call, gotArgs...)
		calls = append(calls, call)
		return nil
	}
	t.Cleanup(func() {
		run = previousRun
	})

	err := GitCommitPostStep{log: zerolog.Nop()}.Run(
		context.Background(),
		PostStepInput{ProjectPath: "/tmp/project"},
	)
	require.NoError(t, err)

	assert.Equal(
		t,
		[][]string{
			{"/tmp/project", "git", "add", "."},
			{"/tmp/project", "git", "commit", "-m", "Initial commit"},
		},
		calls,
	)
}

func TestGitInitPostStepStopsOnError(t *testing.T) {
	expectedErr := errors.New("boom")
	var calls int
	previousRun := run
	run = func(_ context.Context, _ zerolog.Logger, _ string, _ string, _ ...string) error {
		calls++
		return expectedErr
	}
	t.Cleanup(func() {
		run = previousRun
	})

	err := GitInitPostStep{log: zerolog.Nop()}.Run(
		context.Background(),
		PostStepInput{ProjectPath: "/tmp/project"},
	)
	require.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, calls)
}

func TestGitCommitPostStepStopsOnError(t *testing.T) {
	expectedErr := errors.New("boom")
	var calls int
	previousRun := run
	run = func(_ context.Context, _ zerolog.Logger, _ string, _ string, _ ...string) error {
		calls++
		return expectedErr
	}
	t.Cleanup(func() {
		run = previousRun
	})

	err := GitCommitPostStep{log: zerolog.Nop()}.Run(
		context.Background(),
		PostStepInput{ProjectPath: "/tmp/project"},
	)
	require.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, calls)
}
