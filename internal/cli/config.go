package cli

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var cfg config

type config struct {
	Types []Type `toml:"type"`
}

func (c config) toMap() map[string]Type {
	out := make(map[string]Type, len(c.Types))

	for _, t := range c.Types {
		out[t.Name] = t
	}

	return out
}

// Type - A GraphQL type in a GraphQL schema.
type Type struct {
	// Name - The name of the GraphQL type you want to keep/pick.
	Name string `toml:"name"`

	// NewName - The name we should use when creating the new type.
	//
	// Optional.
	//
	// This field exists for scenarios where the original/source type cannot be
	// removed from the existing GraphQL schema, and you need to come up with a
	// new name for your type variant to avoid a name collision.
	NewName string `toml:"new_name"`

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
