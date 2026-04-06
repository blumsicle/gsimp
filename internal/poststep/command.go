package poststep

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog"
)

var run = runCommand

func runCommand(
	ctx context.Context,
	log zerolog.Logger,
	dir string,
	name string,
	args ...string,
) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	log.Debug().
		Str("dir", dir).
		Str("command", name).
		Strs("args", args).
		Msg("executing command")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug().
			Str("dir", dir).
			Str("command", name).
			Strs("args", args).
			Bytes("output", output).
			Msg("command failed")
		if len(output) == 0 {
			return err
		}
		return fmt.Errorf("%w: %s", err, output)
	}

	log.Debug().
		Str("dir", dir).
		Str("command", name).
		Strs("args", args).
		Msg("command completed")

	return nil
}
