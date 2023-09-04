package cli

import (
	"bufio"
	"github.com/charmbracelet/log"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
	"os"
)

func printSDL(doc *ast.Document) {
	sdl := printer.Print(doc)

	// Create a new file where we'll write the new SDL
	// TODO make configurable
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
