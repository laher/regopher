package main

import (
	"fmt"

	"github.com/dave/dst"
)

const cmdParamsToStruct = "params-to-struct"

func regopherParamsToStruct(p inputPos, files map[string]*dst.File, fn *dst.FuncDecl) (map[string]*dst.File, error) {
	structName := fn.Name.Name + "Param"
	varName := "param"
	params := dst.Clone(fn.Type.Params).(*dst.FieldList)
	f := files[p.file]

	fn.Type.Params.List = []*dst.Field{
		structAsField(varName, structName),
	}
	type rename struct {
		from  string
		to    string
		param *dst.Field
	}
	renames := map[string]rename{}
	// rename all references to the parameter
	for _, p := range params.List {
		for _, n := range p.Names {
			renames[n.Name] = rename{
				from:  n.Name,
				to:    fmt.Sprintf("%s.%s", varName, n.Name),
				param: p,
			}
		}
	}
	dst.Inspect(fn.Body, func(n dst.Node) bool {
		switch i := n.(type) {
		case *dst.Ident:
			if n, ok := renames[i.Name]; ok {
				//fmt.Printf("%+v: %+v\n", n.param.Type, i.Obj.Decl.(*dst.Field).Type)
				pt, ok := n.param.Type.(*dst.Ident)
				f, ok := i.Obj.Decl.(*dst.Field)
				if ok {
					nodet, ok := f.Type.(*dst.Ident)
					if ok {
						if pt.String() == nodet.String() {
							i.Name = n.to
						}
					}
				}
			}
		}
		return true
	})
	decl := newStruct(structName, params)
	f.Decls = append(f.Decls, decl)
	updated := map[string]*dst.File{p.file: f}
	for filename, f := range files {
		dst.Inspect(f, func(n dst.Node) bool {
			switch c := n.(type) {
			case *dst.CallExpr:
				switch ce := c.Fun.(type) {
				case *dst.Ident:
					// refactor params
					if ce.Name == fn.Name.Name {
						// TODO: keyify?
						args := c.Args
						c.Args = []dst.Expr{
							&dst.CompositeLit{
								Type: &dst.Ident{
									Name: structName,
								},
								Elts: args,
							},
						}
						updated[filename] = f
					}
				case *dst.SelectorExpr:
					// TODO: selector-based references (perhaps only needed for references
					if ce.Sel.Name == fn.Name.Name {
						// TODO
						fmt.Println("TODO: selector calls not implemented *yet*")
					}
				default:
					// TODO any other types of reference?
					//fmt.Printf("[%s] call to func: %T - %+v\n", filename, ce, ce)
				}
			}
			return true
		})
	}
	// TODO other packages
	return updated, nil
}
