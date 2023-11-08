package cli

import (
	"bufio"
	"os"

	"github.com/charmbracelet/log"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
)

func printSDL(doc *ast.Document) {
	log.Infof("Printing new SDL to file with %d definitions", len(doc.Definitions))

	// Create a new file where we'll write the new SDL
	f, err := os.Create(output)
	if err != nil {
		log.Fatal("unable to create new SDL file for writing", "err", err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	sdl := printer.Print(doc)

	sdlString, ok := sdl.(string)
	if !ok {
		log.Fatal("expected SDL to be a string")
	}

	log.Debug(sdlString)

	_, err = w.WriteString(sdlString)
	if err != nil {
		log.Fatal("unable to produce new SDL file", "err", err)
	}
	defer w.Flush()
}
