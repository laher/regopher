package main

import (
	"fmt"

	"github.com/dave/dst"
)

func regopherParamsToStruct(p string, files map[string]*dst.File, fn *dst.FuncDecl) (map[string]*dst.File, error) {
	paramName := fn.Name.Name + "Param"
	varName := "param"
	params := dst.Clone(fn.Type.Params).(*dst.FieldList)
	f := files[p.file]

	fn.Type.Params.List = []*dst.Field{
		//structAsField(varName, paramName),
	}
	renames := map[string]string{}
	// rename all references to the parameter
	for _, pa := range params.List {
		for _, n := range pa.Names {
			renames[n.Name] = fmt.Sprintf("%s.%s", varName, n.Name)
		}
	}
	return nil, nil
}
