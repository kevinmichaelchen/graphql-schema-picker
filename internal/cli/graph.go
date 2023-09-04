package cli

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/dominikbraun/graph"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/sanity-io/litter"
)

type Vertex struct {
	Name string
	Node ast.Node
}

func NewVertex(node ast.Node) Vertex {
	var name string
	switch node.GetKind() {
	case kinds.ScalarDefinition:
		obj := node.(*ast.ScalarDefinition)
		name = obj.GetName().Value
	case kinds.InterfaceDefinition:
		obj := node.(*ast.InterfaceDefinition)
		name = obj.GetName().Value
	case kinds.UnionDefinition:
		obj := node.(*ast.UnionDefinition)
		name = obj.GetName().Value
	case kinds.EnumDefinition:
		obj := node.(*ast.EnumDefinition)
		name = obj.GetName().Value
	case kinds.InputObjectDefinition:
		obj := node.(*ast.InputObjectDefinition)
		name = obj.GetName().Value
	case kinds.ObjectDefinition:
		obj := node.(*ast.ObjectDefinition)
		name = obj.GetName().Value
	default:
		panic("NewVertex: unsupported node kind: " + node.GetKind())
	}

	return Vertex{
		Name: name,
		Node: node,
	}
}

func VertexHash(v Vertex) string {
	return v.Name
}

// Constructs a graph from the schema AST.
//
// This will provide constant-time access to a map of type names to their
// definitions, as well as definitions to their dependencies (other types).
func buildSchemaGraph(doc *ast.Document) graph.Graph[string, Vertex] {
	g := graph.New(VertexHash)

	// Load the important definitions:
	// scalars, objects, inputs, interfaces, unions, enums
	loadTopLevelDefinitions(g, doc)

	// Build edges between vertices.
	buildEdges(g, doc)

	return g
}

func loadTopLevelDefinitions(g graph.Graph[string, Vertex], doc *ast.Document) {
	for i, d := range doc.Definitions {
		switch d.GetKind() {
		case kinds.ScalarDefinition:
			fallthrough
		case kinds.ObjectDefinition:
			fallthrough
		case kinds.InterfaceDefinition:
			fallthrough
		case kinds.UnionDefinition:
			fallthrough
		case kinds.EnumDefinition:
			fallthrough
		case kinds.InputObjectDefinition:
			v := NewVertex(d)
			_ = g.AddVertex(v)
			log.Debugf("Adding vertex for definition %d (%s) -- %s", i, d.GetKind(), v.Name)

		default:
			log.Warnf("Ignoring definition %d (%s)", i, d.GetKind())
		}
	}
}

func buildEdges(g graph.Graph[string, Vertex], doc *ast.Document) {
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
			obj := v.Node.(*ast.ObjectDefinition)
			fields := obj.Fields
			// TODO iterate through node's fields
			for _, f := range fields {
				rootType := getRootTypeNameHelper(f.Type, 0)
				if isBasicType(rootType) {
					continue
				}
				log.Debug("field",
					"name", f.Name.Value,
					"type", rootType,
				)
				litter.Dump(f.Arguments)
			}
		}

	}

}

func isBasicType(t string) bool {
	return t == "String" || t == "Float" || t == "Int" || t == "Boolean" || t == "ID"
}

// This function returns the Name of the GraphQL type.
//
// There are 3 kinds of "Types" in GraphQL:
//   - Named
//   - List
//   - NonNull
//
// The latter two are wrappers of the Name and therefore require recursion.
func getRootTypeNameHelper(t ast.Type, recursionCount int) string {
	if v, ok := t.(*ast.Named); ok {
		return v.Name.Value
	}

	if v, ok := t.(*ast.List); ok {
		return getRootTypeNameHelper(v.Type, recursionCount+1)
	}

	if v, ok := t.(*ast.NonNull); ok {
		return getRootTypeNameHelper(v.Type, recursionCount+1)
	}

	panic("invalid *ast.Type")
}
