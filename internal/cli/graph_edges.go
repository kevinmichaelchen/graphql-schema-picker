package cli

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/dominikbraun/graph"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func buildEdges(g graph.Graph[string, Vertex]) {
	for _, desired := range desiredDefinitions {
		v, err := g.Vertex(desired)
		if err != nil {
			if errors.Is(err, graph.ErrVertexNotFound) {
				log.Errorf("unable to find definition for: %s", desired)
			} else {
				log.Fatal("unable to read vertex", "err", err)
			}
		}

		switch v.Node.GetKind() {
		case kinds.ObjectDefinition:
			buildEdgesForObject(v, g)
		case kinds.InterfaceDefinition:
			buildEdgesForInterface(v, g)
		case kinds.UnionDefinition:
			buildEdgesForUnion(v, g)
		case kinds.InputObjectDefinition:
			buildEdgesForInputObject(v, g)
		default:
			log.Warnf("Ignoring dependencies for %s node", v.Node.GetKind())
		}
	}
}

func buildEdgesFromFieldDefs(
	g graph.Graph[string, Vertex],
	name string,
	fields []*ast.FieldDefinition,
) {
	for _, f := range fields {
		// Is the field a primitive scalar (e.g., Int, String)?
		// If so, we can skip it, as it's natively a part of any
		// GraphQL schema.
		rootType := getRootTypeNameHelper(f.Type, 0)
		if isBasicType(rootType) {
			continue
		}

		log.Debug("Found field in object",
			"object", name,
			"name", f.Name.Value,
			"type", rootType,
		)

		_ = g.AddEdge(name, rootType)

		// Iterate through f.Argument --
		// since Fields are also dependencies themselves!
		args := f.Arguments
		for _, arg := range args {
			root := getRootTypeNameHelper(arg.Type, 0)
			if isBasicType(root) {
				continue
			}
			_ = g.AddEdge(name, root)
		}
	}
}

func buildEdgesForObject(v Vertex, g graph.Graph[string, Vertex]) {
	obj := v.Node.(*ast.ObjectDefinition)
	fields := obj.Fields
	buildEdgesFromFieldDefs(g, obj.Name.Value, fields)
}

func buildEdgesForInterface(v Vertex, g graph.Graph[string, Vertex]) {
	obj := v.Node.(*ast.InterfaceDefinition)
	fields := obj.Fields
	buildEdgesFromFieldDefs(g, obj.Name.Value, fields)
}

func buildEdgesForUnion(v Vertex, g graph.Graph[string, Vertex]) {
	// TODO add support
}

func buildEdgesForInputObject(v Vertex, g graph.Graph[string, Vertex]) {
	obj := v.Node.(*ast.InputObjectDefinition)
	fields := obj.Fields
	name := obj.Name.Value

	for _, f := range fields {
		// Is the field a primitive scalar (e.g., Int, String)?
		// If so, we can skip it, as it's natively a part of any
		// GraphQL schema.
		rootType := getRootTypeNameHelper(f.Type, 0)
		if isBasicType(rootType) {
			continue
		}

		log.Debug("Found field in object",
			"object", name,
			"name", f.Name.Value,
			"type", rootType,
		)

		_ = g.AddEdge(name, rootType)
	}
}
