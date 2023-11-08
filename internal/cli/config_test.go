package cli

import (
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"
	"testing"
)

const blob = `
[[type]]
name = "PersonInsertInput"
new_name = "SvcPersonInsertInput"
deny_list = ["address"]

[[type]]
name = "Person"
new_name = "SvcPerson"
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
				NewName:  "SvcPersonInsertInput",
				DenyList: []string{"address"},
			},
			{
				Name:     "Person",
				NewName:  "SvcPerson",
				DenyList: []string{"address"},
			},
		},
	}

	require.Equal(t, expected, testCfg)
}
