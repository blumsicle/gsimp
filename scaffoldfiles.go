// Package bcli embeds canonical source files used by project generation.
package bcli

import "embed"

// ScaffoldSourceFS contains canonical source files copied into generated projects.
//
//go:embed internal/cli/runner.go internal/cli/buildinfo.go internal/appconfig/load.go
var ScaffoldSourceFS embed.FS
