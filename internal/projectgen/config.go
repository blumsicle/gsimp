package projectgen

import "fmt"

// Config defines the inputs used to generate a new project scaffold.
type Config struct {
	Name             string
	Description      string
	GitLocation      string
	ProjectDirPrefix string
	RootPath         string
	InPlace          bool
}

func validateConfig(cfg Config) error {
	if cfg.Name == "" {
		return fmt.Errorf("name is required")
	}
	if cfg.Description == "" {
		return fmt.Errorf("description is required")
	}

	return nil
}
