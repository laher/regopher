package main

import (
	"fmt"

	"github.com/dave/dst"
)

func extractParameterObject(f *dst.File, funcName string) error {
	var params *dst.FieldList
	paramName := funcName + "Param"
	for _, d := range f.Decls {
		//	fmt.Printf("%T: %+v\n", d, d)
		switch fn := d.(type) {
		case *dst.FuncDecl:
			if fn.Name.Name == funcName {

				params = dst.Clone(fn.Type.Params).(*dst.FieldList)

				fn.Type.Params.List = []*dst.Field{
					structAsField("param", paramName),
				}
			}
		}
	}
	if params == nil {
		return fmt.Errorf("function '%s' not found", funcName)
	}
	decl := newStruct(paramName, params)
	f.Decls = append(f.Decls, decl)

	return nil
}
