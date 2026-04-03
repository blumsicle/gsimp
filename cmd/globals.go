package cmd

type Globals struct {
	CredsFile string `short:"c" default:"~/.bcreds.json" type:"path" help:"Credentials file to use"`
}
