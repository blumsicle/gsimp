package poststep

import (
	"context"
	"errors"
	"testing"

	"github.com/blumsicle/bcli/internal/bcliconfig"
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

func TestCommandPostStepSpecReturnsRunError(t *testing.T) {
	expectedErr := errors.New("boom")
	previousRun := run
	run = func(_ context.Context, _ zerolog.Logger, _ string, _ string, _ ...string) error {
		return expectedErr
	}
	t.Cleanup(func() {
		run = previousRun
	})

	err := commandPostStepSpec{
		name:    "test",
		message: "testing command",
		command: "test",
		args:    []string{"arg"},
	}.run(context.Background(), zerolog.Nop(), PostStepInput{ProjectPath: "/tmp/project"})

	require.ErrorIs(t, err, expectedErr)
}

func TestPlannedReturnsDefaultStepsInOrder(t *testing.T) {
	cfg := bcliconfig.Default()

	steps := NewPlanner(zerolog.Nop(), &cfg.PostSteps).Planned()

	assert.Equal(
		t,
		[]string{"go get -u ./...", "go mod tidy", "git init", "git commit"},
		postStepNames(steps),
	)
}

func TestPlannedSkipsDisabledSteps(t *testing.T) {
	tests := []struct {
		name      string
		configure func(cfg *bcliconfig.PostStepsConfig)
		want      []string
	}{
		{
			name: "go get update",
			configure: func(cfg *bcliconfig.PostStepsConfig) {
				cfg.GoGetUpdate = false
			},
			want: []string{"go mod tidy", "git init", "git commit"},
		},
		{
			name: "go mod tidy",
			configure: func(cfg *bcliconfig.PostStepsConfig) {
				cfg.GoModTidy = false
			},
			want: []string{"go get -u ./...", "git init", "git commit"},
		},
		{
			name: "git init",
			configure: func(cfg *bcliconfig.PostStepsConfig) {
				cfg.GitInit = false
			},
			want: []string{"go get -u ./...", "go mod tidy"},
		},
		{
			name: "git commit",
			configure: func(cfg *bcliconfig.PostStepsConfig) {
				cfg.GitCommit = false
			},
			want: []string{"go get -u ./...", "go mod tidy", "git init"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := bcliconfig.Default()
			tt.configure(&cfg.PostSteps)

			steps := NewPlanner(zerolog.Nop(), &cfg.PostSteps).Planned()

			assert.Equal(t, tt.want, postStepNames(steps))
		})
	}
}

func TestPlannedDisablesGitCommitWhenGitInitIsDisabled(t *testing.T) {
	cfg := bcliconfig.Default()
	cfg.PostSteps.GitInit = false
	cfg.PostSteps.GitCommit = true

	steps := NewPlanner(zerolog.Nop(), &cfg.PostSteps).Planned()

	assert.Equal(
		t,
		[]string{"go get -u ./...", "go mod tidy"},
		postStepNames(steps),
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

func postStepNames(steps []PostStep) []string {
	names := make([]string, 0, len(steps))
	for _, step := range steps {
		names = append(names, step.Name())
	}

	return names
}
