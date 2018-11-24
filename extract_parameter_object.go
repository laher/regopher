package main

import (
	"fmt"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const introduceParameterObject = "introduce-parameter-object"

func getIdentifierAt(d *decorator.Decorator, f *dst.File, pos int) (*dst.Ident, error) {
	var r *dst.Ident
	var err error
	dst.Inspect(f, func(n dst.Node) bool {
		if n == nil {
			return false
		}
		switch n.(type) {
		case *dst.FuncType: //ignore.
			// TODO: Why FuncType, specifically?
			// Ignore any others?
			return true
		}

		an, ok := d.Ast.Nodes[n]
		if !ok {
			// TODO: ignore like FuncType?
			err = fmt.Errorf("No ast node matching dst for %T: %+v", n, n)
			return false
		}
		p := an.Pos()
		e := an.End()
		if pos > int(p) && pos < int(e) {
			if i, ok := n.(*dst.Ident); ok {
				// fmt.Printf("%T [%d,%d]: %+v\n", n, p, e, n)
				r = i
			}
			// ok
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, fmt.Errorf("ident not found")
	}
	return r, nil
}

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
	}
	if params == nil {
		return fmt.Errorf("function '%s' not found", funcName)
	}
	decl := newStruct(paramName, params)
	f.Decls = append(f.Decls, decl)

	return nil
}
