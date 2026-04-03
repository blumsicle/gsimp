package example

import (
	"github.com/blumsicle/gsimp/cmd"
	"github.com/rs/zerolog"
)

type Command struct{}

func (c *Command) Run(log zerolog.Logger, g *cmd.Globals) error {
	log.Info().Str("config_file", g.ConfigFile).Msg("example command")
	return nil
}
