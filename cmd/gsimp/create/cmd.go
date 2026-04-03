package create

import (
	"github.com/blumsicle/gsimp/cmd"
	"github.com/rs/zerolog"
)

type Command struct{}

func (c *Command) Run(log zerolog.Logger, g *cmd.Globals) error {
	log.Info().Msg("in create")
	return nil
}
