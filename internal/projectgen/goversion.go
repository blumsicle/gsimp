package projectgen

import (
	"os/exec"
	"runtime"
	"strings"
)

func currentGoVersion() string {
	version, err := exec.Command("go", "env", "GOVERSION").Output()
	if err == nil {
		goVersion := strings.TrimSpace(string(version))
		goVersion = strings.TrimPrefix(goVersion, "go")
		if goVersion != "" {
			return goVersion
		}
	}

	return strings.TrimPrefix(runtime.Version(), "go")
}
