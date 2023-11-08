package cli

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "graphql-schema-picker",
	Short: "Utility for picking certain definitions from a GraphQL schema into a new schema",
	Long:  `Utility for picking certain definitions from a GraphQL schema (SDL file) into a smaller schema`,
	Run: func(cmd *cobra.Command, args []string) {
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		err := loadConfig()
		if err != nil {
			log.Fatal("error occurred while loading config", "err", err)
		}
	},
}

var (
	ldFlags            LDFlags
	sdlFile            string
	debug              bool
	dryRun             bool
	desiredDefinitions []string
	output             string
	configPath         string
)

// LDFlags contain fields that get linked and compiled into the final binary
// program at build time.
type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func init() {
	rootCmd.AddCommand(pick)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVarP(&sdlFile, "sdl-file", "f", "", "path to an SDL file")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "verbose debug logging")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "", false, "dry run")

	pick.Flags().StringSliceVarP(&desiredDefinitions, "definitions", "d", []string{}, "definitions from the SDL you want to pick/keep")
	pick.Flags().StringVarP(&output, "output", "o", "", "where the resulting schema/SDL file is written")
	pick.Flags().StringVarP(&configPath, "config", "c", "", "path to config file")
}

func Main(ldf LDFlags) {
	ldFlags = ldf
	if err := rootCmd.Execute(); err != nil {
		log.Error("execution failed", "err", err)
		os.Exit(1)
	}
}
