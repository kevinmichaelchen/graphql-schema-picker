package cli

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"
	"testing"
)

const blob = `
[[type]]
name = "Person"
deny_list = ["foobar"]

[[type]]
name = "Stairway to Heaven"
`

func TestDecode(t *testing.T) {
	var testCfg config

	_, err := toml.Decode(blob, &testCfg)
	require.NoError(t, err)

	for _, s := range testCfg.Types {
		fmt.Printf("%s (%v)\n", s.Name, s.DenyList)
	}
}
