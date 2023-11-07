package cli

import (
	"github.com/charmbracelet/log"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

func parseDocument(source string) *ast.Document {
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: source,
		Options: parser.ParseOptions{
			NoLocation: true,
			NoSource:   true,
		},
	})
	if err != nil {
		log.Fatal("unable to parse SDL", "err", err)
	}

	log.Infof("Parsed document with %d definitions", len(astDoc.Definitions))

	return astDoc
}
