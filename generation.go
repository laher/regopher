package main

import (
	"go/token"

	"github.com/dave/dst"
)

func structAsField(varName, typeName string) *dst.Field {
	return &dst.Field{
		Names: []*dst.Ident{
			&dst.Ident{Name: varName},
		},
		Type: &dst.Ident{Name: typeName},
	}
}

func newStruct(name string, fields *dst.FieldList) *dst.GenDecl {
	decl := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{Name: name},
				Type: &dst.StructType{
					Fields: fields,
				},
			},
		},
		Decs: dst.GenDeclDecorations{
			NodeDecs: dst.NodeDecs{},
		},
	}
	return decl
}
