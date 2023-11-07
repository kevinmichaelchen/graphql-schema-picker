package cli

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/dominikbraun/graph/draw"
	"github.com/spf13/cobra"
)

var pick = &cobra.Command{
	Use:   "pick",
	Short: "Picks certain definitions from a GraphQL Schema (SDL) file",
	Long:  `Picks certain definitions from a GraphQL Schema (SDL) file`,
	Run:   fnPick,
}

func fnPick(cmd *cobra.Command, args []string) {
	if len(desiredDefinitions) == 0 {
		log.Fatal("you must provide at least one definition you wish to pick")
	}

	// Open the file
	b, err := os.ReadFile(sdlFile)
	if err != nil {
		log.Fatal("unable to read file", "err", err)
	}

	// Parse it into an Abstract Syntax Tree (AST)
	astDoc := parseDocument(string(b))

	// The most crucial step is to map out the AST.
	// The AST object itself is clunky and nested.
	// We need something flatter and more ergonomic: like a couple of maps!
	g := buildPrunedGraph(astDoc)

	// Build a visual diagram of the new SDL
	file, _ := os.Create("./simple.gv")
	_ = draw.DOT(g, file)

	// Convert the GRAPH back to an SDL
	newSDL := convertGraphToSDL(g)

	// Print the new SDL to file
	printSDL(newSDL)
}
