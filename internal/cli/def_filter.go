package cli

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

// filters fields out of a GraphQL type definition
func filterDef(def Vertex) ast.Node {
	var out ast.Node

	node := def.Node

	switch node.GetKind() {
	case kinds.ObjectDefinition:
		obj := node.(*ast.ObjectDefinition)
		objName := obj.Name.Value

		newObjName := getNewObjectName(objName)

		var fields []*ast.FieldDefinition
		for _, field := range obj.Fields {
			fieldName := field.Name.Value

			if isFilteredField(objName, fieldName) {
				continue
			}

			fields = append(fields, field)
		}

		out = &ast.ObjectDefinition{
			Kind: obj.Kind,
			// Location might not be accurate since we're pruning, so exclude it
			//Loc: obj.Loc,
			Name: &ast.Name{
				Kind:  kinds.StringValue,
				Loc:   nil,
				Value: newObjName,
			},
			Description: obj.Description,
			// TODO we don't support interface pruning -
			//  if we prune a field on the interface, we'd have to prune the
			//  interface as well - not sure how to handle that yet
			Interfaces: nil,
			Directives: obj.Directives,
			Fields:     fields,
		}

	case kinds.InputObjectDefinition:
		obj := node.(*ast.InputObjectDefinition)
		objName := obj.Name.Value

		newObjName := getNewObjectName(objName)

		var fields []*ast.InputValueDefinition
		for _, field := range obj.Fields {
			fieldName := field.Name.Value

			if isFilteredField(objName, fieldName) {
				continue
			}

			fields = append(fields, field)
		}

		out = &ast.InputObjectDefinition{
			Kind: obj.Kind,
			// Location might not be accurate since we're pruning, so exclude it
			//Loc: obj.Loc,
			Name: &ast.Name{
				Kind:  kinds.StringValue,
				Loc:   nil,
				Value: newObjName,
			},
			Description: obj.Description,
			Directives:  obj.Directives,
			Fields:      fields,
		}

	default:
		return node
	}

	return out
}

func getNewObjectName(objName string) string {
	defCfg := cfg.toMap()[objName]
	if defCfg.NewName != "" {
		return defCfg.NewName
	}

	return defCfg.Name
}

// isFilteredField - Checks if the field should be filtered out
func isFilteredField(objName, fieldName string) bool {
	defCfg := cfg.toMap()[objName]
	denyList := defCfg.DenyList

	for _, e := range denyList {
		if fieldName == e {
			return true
		}
	}

	return false
}
