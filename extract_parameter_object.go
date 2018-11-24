package main

import (
	"fmt"

	"github.com/dave/dst"
)

const introduceParameterObject = "introduce-parameter-object"

func extractParameterObject(p inputPos, files map[string]*dst.File, fn *dst.FuncDecl) (map[string]*dst.File, error) {
	paramName := fn.Name.Name + "Param"
	varName := "param"
	params := dst.Clone(fn.Type.Params).(*dst.FieldList)
	f := files[p.file]

	fn.Type.Params.List = []*dst.Field{
		structAsField(varName, paramName),
	}
	renames := map[string]string{}
	// rename all references to the parameter
	for _, p := range params.List {
		for _, n := range p.Names {
			renames[n.Name] = fmt.Sprintf("%s.%s", varName, n.Name)
		}
	}
	dst.Inspect(fn.Body, func(n dst.Node) bool {
		switch i := n.(type) {
		case *dst.Ident:
			if n, ok := renames[i.Name]; ok {
				i.Name = n
			}
		}
		return true
	})
	decl := newStruct(paramName, params)
	f.Decls = append(f.Decls, decl)

	// TODO other files
	return map[string]*dst.File{p.file: f}, nil
}
