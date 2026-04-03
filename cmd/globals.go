package cmd

type Globals struct {
	ConfigFile string `short:"c" default:"~/.config/starter/config.yaml" type:"path" help:"Path to the config file"`
}
