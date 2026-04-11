// Package cli contains shared runtime helpers for CLI binaries.
package cli

import "runtime/debug"

const shortCommitLength = 7

const (
	defaultVersion = "dev"
	defaultCommit  = "unknown"
)

// BuildInfo describes resolved build metadata for a CLI binary.
type BuildInfo struct {
	Name    string
	Version string
	Commit  string
}

// ResolveBuildInfo returns build metadata with runtime build info fallbacks.
func ResolveBuildInfo(name string) BuildInfo {
	info := BuildInfo{
		Name:    name,
		Version: defaultVersion,
		Commit:  defaultCommit,
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}

	return applyRuntimeBuildInfo(info, *buildInfo)
}

func applyRuntimeBuildInfo(info BuildInfo, buildInfo debug.BuildInfo) BuildInfo {
	if (info.Version == "" || info.Version == defaultVersion) &&
		buildInfo.Main.Version != "" &&
		buildInfo.Main.Version != "(devel)" {
		info.Version = buildInfo.Main.Version
	}

	if info.Commit == "" || info.Commit == defaultCommit {
		if revision := buildSetting(buildInfo.Settings, "vcs.revision"); revision != "" {
			info.Commit = shortCommit(revision)
		}
	}

	return info
}

func shortCommit(revision string) string {
	if len(revision) <= shortCommitLength {
		return revision
	}

	return revision[:shortCommitLength]
}

func buildSetting(settings []debug.BuildSetting, key string) string {
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}

	return ""
}
