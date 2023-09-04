package cli

import (
	"github.com/charmbracelet/log"
	"github.com/dominikbraun/graph"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func convertGraphToSDL(g graph.Graph[string, Vertex]) *ast.Document {
	out := &ast.Document{
		Kind:        kinds.Document,
		Loc:         nil,
		Definitions: nil,
	}

	var defs []ast.Node

	adj, err := g.AdjacencyMap()
	if err != nil {
		log.Fatal("unable to retrieve AdjacencyMap", "err", err)
	}

	for defName := range adj {
		def := must(g.Vertex(defName))
		defs = append(defs, def.Node)
	}

	out.Definitions = defs

	return out
}
