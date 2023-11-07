package cli

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/dominikbraun/graph"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

type Vertex struct {
	Name string
	Node ast.Node
}

type Describer interface {
	GetDescription() *ast.StringValue
}

func sanitizeComment(d Describer) string {
	desc := d.GetDescription()
	if desc == nil {
		return ""
	}

	return strings.ReplaceAll(desc.Value, `"`, `'`)
}

func NewVertex(node ast.Node) Vertex {
	var name string

	switch node.GetKind() {
	case kinds.ScalarDefinition:
		obj := node.(*ast.ScalarDefinition)
		name = obj.GetName().Value

		// Sanitize description (e.g., remove double-quotes)
		if obj.Description != nil {
			obj.Description = &ast.StringValue{
				Kind:  kinds.StringValue,
				Value: sanitizeComment(obj),
			}
		}
	case kinds.InterfaceDefinition:
		obj := node.(*ast.InterfaceDefinition)
		name = obj.GetName().Value

		// Sanitize description (e.g., remove double-quotes)
		if obj.Description != nil {
			obj.Description = &ast.StringValue{
				Kind:  kinds.StringValue,
				Value: sanitizeComment(obj),
			}
		}
	case kinds.UnionDefinition:
		obj := node.(*ast.UnionDefinition)
		name = obj.GetName().Value

		// Sanitize description (e.g., remove double-quotes)
		if obj.Description != nil {
			obj.Description = &ast.StringValue{
				Kind:  kinds.StringValue,
				Value: sanitizeComment(obj),
			}
		}
	case kinds.EnumDefinition:
		obj := node.(*ast.EnumDefinition)
		name = obj.GetName().Value

		// Sanitize description (e.g., remove double-quotes)
		if obj.Description != nil {
			obj.Description = &ast.StringValue{
				Kind:  kinds.StringValue,
				Value: sanitizeComment(obj),
			}
		}
	case kinds.InputObjectDefinition:
		obj := node.(*ast.InputObjectDefinition)
		name = obj.GetName().Value

		// Sanitize description (e.g., remove double-quotes)
		if obj.Description != nil {
			obj.Description = &ast.StringValue{
				Kind:  kinds.StringValue,
				Value: sanitizeComment(obj),
			}
		}

	case kinds.ObjectDefinition:
		obj := node.(*ast.ObjectDefinition)
		name = obj.GetName().Value

		// Sanitize description (e.g., remove double-quotes)
		if obj.Description != nil {
			obj.Description = &ast.StringValue{
				Kind:  kinds.StringValue,
				Value: sanitizeComment(obj),
			}
		}
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
func buildPrunedGraph(doc *ast.Document) graph.Graph[string, Vertex] {
	g := graph.New(VertexHash)

	// Load the important definitions:
	// scalars, objects, inputs, interfaces, unions, enums
	loadTopLevelDefinitions(g, doc)

	// Build edges between vertices.
	buildEdges(g)

	// Prune our graph:
	// 1. We filter out any vertices (GQL types) that we don't explicitly want
	// 2. We filter out any fields within those GQL types
	return prune(g)
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
			log.Debugf("Adding vertex for definition %d (%s) -- %s",
				i, d.GetKind(), v.Name)

		default:
			log.Warnf("Ignoring definition %d (%s)", i, d.GetKind())
		}
	}
}

// Prune our graph:
// 1. We filter out any vertices (GQL types) that we don't explicitly want
// 2. We filter out any fields within those GQL types
func prune(in graph.Graph[string, Vertex]) graph.Graph[string, Vertex] {
	m, err := in.AdjacencyMap()
	if err != nil {
		log.Fatal("unable to retrieve adjacency map", "err", err)
	}

	out := graph.New(VertexHash)
	for defName, edges := range m {
		if !isDesired(defName) {
			continue
		}

		// Retrieve the vertex (GraphQL type) by its name
		def := must(in.Vertex(defName))

		// TODO clone the type definition and filter out fields
		//  def = filter(clone(def))

		// Add vertex (GraphQL type) to our new, outgoing graph
		_ = out.AddVertex(def)

		// Add its dependent vertices
		for depName := range edges {
			dep := must(in.Vertex(depName))
			_ = out.AddVertex(dep)
			_ = out.AddEdge(defName, depName)
		}
	}

	return out
}

func must(v Vertex, err error) Vertex {
	if err != nil {
		log.Fatal("unable to find vertex", "err", err)
	}
	return v
}

func isDesired(definitionName string) bool {
	// TODO deprecate in favor of config
	for _, d := range desiredDefinitions {
		if definitionName == d {
			return true
		}
	}

	return false
}

func isBasicType(t string) bool {
	// https://graphql.org/graphql-js/basic-types/
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
