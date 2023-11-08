package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints CLI version info",
	Long:  `Prints CLI version info`,
	Run:   fnVersion,
}

func fnVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("version %s\n", ldFlags.Version)
}
