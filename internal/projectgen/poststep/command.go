package poststep

import (
	"context"
	"fmt"
	"os/exec"
)

func runCommand(ctx context.Context, dir string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) == 0 {
			return err
		}
		return fmt.Errorf("%w: %s", err, output)
	}

	return nil
}
