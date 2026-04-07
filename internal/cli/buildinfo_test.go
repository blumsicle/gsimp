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
			Version: "dev",
			Commit:  "unknown",
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
	assert.Equal(t, "abcdef123456", info.Commit)
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
			Version: "dev",
			Commit:  "unknown",
		},
		debug.BuildInfo{
			Main: debug.Module{
				Version: "(devel)",
			},
		},
	)

	assert.Equal(t, "dev", info.Version)
	assert.Equal(t, "unknown", info.Commit)
}
