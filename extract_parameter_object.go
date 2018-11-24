package main

import (
	"fmt"

	"github.com/dave/dst"
)

func extractParameterObject(f *dst.File, funcName string) error {
	var params *dst.FieldList
	paramName := funcName + "Param"
	varName := "param"
	for _, d := range f.Decls {
		switch fn := d.(type) {
		case *dst.FuncDecl:
			if fn.Name.Name == funcName {

				params = dst.Clone(fn.Type.Params).(*dst.FieldList)

				fn.Type.Params.List = []*dst.Field{
					structAsField(varName, paramName),
				}
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
		}
	}
	if params == nil {
		return fmt.Errorf("function '%s' not found", funcName)
	}
	decl := newStruct(paramName, params)
	f.Decls = append(f.Decls, decl)

	return nil
}
