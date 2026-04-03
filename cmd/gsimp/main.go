package main

import "github.com/blumsicle/gsimp/cmd"

var (
	name    = "gsimp"
	version = "dev"
	commit  = "unknown"
)

func main() {
	cmd.Run(&CLI{}, "Credential manager", cmd.BuildInfo{
		Name:    name,
		Version: version,
		Commit:  commit,
	})
}
