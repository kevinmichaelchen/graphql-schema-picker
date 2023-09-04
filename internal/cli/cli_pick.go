package cli

import (
	"bufio"
	"github.com/charmbracelet/log"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/spf13/cobra"
	"os"
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
	graph := buildSchemaGraph(astDoc)

	result := graph.getPicked()

	//litter.Dump(result)

	newDoc := &ast.Document{
		Kind:        kinds.Document,
		Loc:         nil,
		Definitions: nil,
	}
	for _, node := range result {
		newDoc.Definitions = append(newDoc.Definitions, node)
	}

	// TODO write new AST to file as SDL
	log.Infof("new doc has %d defs, old one has %d", len(newDoc.Definitions), len(astDoc.Definitions))
	sdl := printer.Print(newDoc)

	//litter.Dump(sdl)
	//log.Infof("sdl is type: %T", sdl)

	f, err := os.Create("output.sdl.graphqls")
	if err != nil {
		log.Fatal("unable to create new SDL file for writing", "err", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(sdl.(string))
	if err != nil {
		log.Fatal("unable to produce new SDL file", "err", err)
	}
}
