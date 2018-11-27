package main

import (
	"fmt"

	"github.com/dave/dst"
)

func regopherParamsToStruct(p string, files map[string]*dst.File, fn *dst.FuncDecl) (map[string]*dst.File, error) {
	varName := "param"
	params := dst.Clone(fn.Type.Params).(*dst.FieldList)
	fn.Type.Params.List = []*dst.Field{
		//structAsField(varName, paramName),
	}
	renames := map[string]string{}
	// rename all references to the parameter
	for _, p := range params.List {
		for _, n := range p.Names {
			renames[n.Name] = fmt.Sprintf("%s.%s", varName, n.Name)
		}
	}
	return nil, nil
}
