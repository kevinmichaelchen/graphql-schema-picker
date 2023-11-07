package cli

import (
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"
	"testing"
)

const blob = `
[[type]]
name = "PersonInsertInput"
deny_list = ["address"]

[[type]]
name = "Person"
deny_list = ["address"]
`

func TestDecode(t *testing.T) {
	var testCfg config

	_, err := toml.Decode(blob, &testCfg)
	require.NoError(t, err)

	expected := config{
		Types: []Type{
			{
				Name:     "PersonInsertInput",
				DenyList: []string{"address"},
			},
			{
				Name:     "Person",
				DenyList: []string{"address"},
			},
		},
	}

	require.Equal(t, expected, testCfg)
}
