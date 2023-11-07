package cli

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var cfg config

type config struct {
	Types []Type `toml:"type"`
}

// Type - A GraphQL type in a GraphQL schema.
type Type struct {
	// Name - The name of the GraphQL type you want to keep/pick.
	Name string `toml:"name"`

	// DenyList - List of fields within that type you want to filter out.
	DenyList []string `toml:"deny_list"`
}

func loadConfig() error {
	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}

	return nil
}
