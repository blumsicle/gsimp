package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionFlag(t *testing.T) {
	harness := newCLITestHarness(t, nil)

	_, err := harness.parse(t, "--version")
	require.Error(t, err)

	assert.Equal(t, 0, harness.exitCode)
	assert.Equal(t, "bcli test-version test-commit\n", harness.stdoutString())
	assert.Empty(t, harness.stderrString())
}

func TestHelpFlag(t *testing.T) {
	harness := newCLITestHarness(t, nil)

	_, err := harness.parse(t, "--help")
	require.Error(t, err)

	assert.Equal(t, 0, harness.exitCode)
	assert.Contains(t, harness.stdoutString(), "Generate starter Go CLI projects")
	assert.Contains(t, harness.stdoutString(), "--log-level")
	assert.Contains(t, harness.stdoutString(), "completion")
	assert.Contains(t, harness.stdoutString(), "config")
	assert.Contains(t, harness.stdoutString(), "create")
	assert.Empty(t, harness.stderrString())
}
