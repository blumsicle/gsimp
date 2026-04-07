// Package cli contains shared runtime helpers for CLI binaries.
package cli

import "runtime/debug"

// BuildInfo describes linker-injected metadata for a CLI binary.
type BuildInfo struct {
	Name    string
	Version string
	Commit  string
}

// ResolveBuildInfo returns build metadata from ldflags with runtime build info fallbacks.
func ResolveBuildInfo(name string, version string, commit string) BuildInfo {
	info := BuildInfo{
		Name:    name,
		Version: version,
		Commit:  commit,
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}

	return applyRuntimeBuildInfo(info, *buildInfo)
}

func applyRuntimeBuildInfo(info BuildInfo, buildInfo debug.BuildInfo) BuildInfo {
	if (info.Version == "" || info.Version == "dev") &&
		buildInfo.Main.Version != "" &&
		buildInfo.Main.Version != "(devel)" {
		info.Version = buildInfo.Main.Version
	}

	if info.Commit == "" || info.Commit == "unknown" {
		if revision := buildSetting(buildInfo.Settings, "vcs.revision"); revision != "" {
			info.Commit = revision
		}
	}

	return info
}

func buildSetting(settings []debug.BuildSetting, key string) string {
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}

	return ""
}
