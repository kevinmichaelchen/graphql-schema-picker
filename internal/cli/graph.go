package cli

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/sanity-io/litter"
)

type schemaGraph struct {
	definitions  map[string]ast.Node
	dependencies map[string]map[string]struct{}
}

func newSchemaGraph() schemaGraph {
	definitions := map[string]ast.Node{}
	deps := map[string]map[string]struct{}{}

	return schemaGraph{
		definitions:  definitions,
		dependencies: deps,
	}
}

func (sg schemaGraph) getPicked() map[string]ast.Node {
	log.Infof("you picked %d definitions: %v", len(desiredDefinitions), desiredDefinitions)

	out := make(map[string]ast.Node, len(desiredDefinitions))
	for _, name := range desiredDefinitions {
		d, ok := sg.definitions[name]
		if !ok {
			log.Fatalf("unable to find definition by name: %s")
		}

		out[name] = d

		// Does this definition depend on anything else?
		// If so, let's add those dependencies to the output as well.
		deps := sg.dependencies[name]
		if len(deps) > 0 {
			log.Info("scanning deps", "def", name, "deps", deps)
			for depName := range deps {
				var depDef ast.Node
				depDef, ok = sg.definitions[depName]
				if !ok {
					log.Fatalf("unable to find dependency definition (for %s) by name: %s", name, depName)
				}

				out[depName] = depDef
			}
		}
	}

	return out
}

func (sg schemaGraph) process(node ast.Node) {
	switch node.GetKind() {
	case kinds.ScalarDefinition:
		s := node.(*ast.ScalarDefinition)
		log.Info("found scalar")
		litter.Dump(s)
		name := s.GetName().Value
		sg.definitions[name] = node

	case kinds.ObjectDefinition:
		obj := node.(*ast.ObjectDefinition)
		name := obj.GetName().Value
		sg.definitions[name] = node
		addDependenciesForObjectDefinition(name, sg.dependencies, obj)
	}
}

// Constructs a graph from the schema AST.
//
// This will provide constant-time access to a map of type names to their
// definitions, as well as definitions to their dependencies (other types).
func buildSchemaGraph(doc *ast.Document) schemaGraph {
	sg := newSchemaGraph()

	for i, d := range doc.Definitions {
		switch d.GetKind() {
		case kinds.ObjectDefinition:
			obj := d.(*ast.ObjectDefinition)
			name := obj.GetName().Value
			log.Infof("Processing definition %d: %s (%s)", i, name, d.GetKind())
		default:
			log.Warnf("Ignoring definition %d (%s)", i, d.GetKind())
		}

		sg.process(d)
	}

	return sg
}

func addDependenciesForObjectDefinition(
	name string,
	deps map[string]map[string]struct{},
	obj *ast.ObjectDefinition,
) {
	for _, fd := range obj.Fields {
		log.Debug("field definition",
			"name.value", fd.Name.Value,
			"kind", fd.GetKind(),
			"type", fd.Type.String(),
		)
		//litter.Dump(fd)

		rootType, err := getRootTypeNameHelper(fd.Type, 0)
		if err != nil {
			log.Fatal("unable to get root type name", "err", err)
		}

		// filter out native scalars (e.g., strings, bools, etc)
		if ignoreDependentType(rootType) {
			continue
		}

		if deps[name] == nil {
			deps[name] = map[string]struct{}{}
		}
		deps[name][rootType] = struct{}{}
	}
}

func ignoreDependentType(t string) bool {
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
func getRootTypeNameHelper(t ast.Type, recursionCount int) (string, error) {
	if v, ok := t.(*ast.Named); ok {
		//log.Infof("found a root type: %s / %s", v.Name.Kind, v.Name.Value)
		if v.Name.Value == "String" {
			litter.Dump(t)
		}
		return v.Name.Value, nil
	}

	if v, ok := t.(*ast.List); ok {
		return getRootTypeNameHelper(v.Type, recursionCount+1)
	}

	if v, ok := t.(*ast.NonNull); ok {
		return getRootTypeNameHelper(v.Type, recursionCount+1)
	}

	return "", errors.New("invalid *ast.Type")
}
