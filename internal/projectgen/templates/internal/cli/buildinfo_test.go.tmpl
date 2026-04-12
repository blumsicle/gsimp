package cli

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyRuntimeBuildInfoUsesRuntimeVersionAndCommitFallbacks(t *testing.T) {
	info := applyRuntimeBuildInfo(
		BuildInfo{
			Name:    "bcli",
			Version: defaultVersion,
			Commit:  defaultCommit,
		},
		debug.BuildInfo{
			Main: debug.Module{
				Version: "v0.2.0",
			},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "abcdef123456"},
			},
		},
	)

	assert.Equal(t, "bcli", info.Name)
	assert.Equal(t, "v0.2.0", info.Version)
	assert.Equal(t, "abcdef1", info.Commit)
}

func TestApplyRuntimeBuildInfoKeepsExplicitLdflagValues(t *testing.T) {
	info := applyRuntimeBuildInfo(
		BuildInfo{
			Name:    "bcli",
			Version: "v9.9.9",
			Commit:  "release-commit",
		},
		debug.BuildInfo{
			Main: debug.Module{
				Version: "v0.2.0",
			},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "abcdef123456"},
			},
		},
	)

	assert.Equal(t, "v9.9.9", info.Version)
	assert.Equal(t, "release-commit", info.Commit)
}

func TestApplyRuntimeBuildInfoIgnoresDevelVersion(t *testing.T) {
	info := applyRuntimeBuildInfo(
		BuildInfo{
			Name:    "bcli",
			Version: defaultVersion,
			Commit:  defaultCommit,
		},
		debug.BuildInfo{
			Main: debug.Module{
				Version: "(devel)",
			},
		},
	)

	assert.Equal(t, defaultVersion, info.Version)
	assert.Equal(t, defaultCommit, info.Commit)
}

func TestApplyRuntimeBuildInfoKeepsShortRevision(t *testing.T) {
	info := applyRuntimeBuildInfo(
		BuildInfo{
			Name:    "bcli",
			Version: defaultVersion,
			Commit:  defaultCommit,
		},
		debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "abc123"},
			},
		},
	)

	assert.Equal(t, "abc123", info.Commit)
}
