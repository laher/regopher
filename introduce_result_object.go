package main

import (
	"fmt"

	"github.com/dave/dst"
)

const cmdIntroduceResultObject = "introduce-result-object"

func introduceResultObject(p inputPos, files map[string]*dst.File, fn *dst.FuncDecl) (map[string]*dst.File, error) {
	resultName := fn.Name.Name + "Result"
	varName := "res"
	results := dst.Clone(fn.Type.Results).(*dst.FieldList)
	f := files[p.file]

	fn.Type.Results.List = []*dst.Field{
		&dst.Field{Type: &dst.Ident{Name: resultName}},
	}

	renames := map[string]string{}
	// rename all references to the result
	for _, p := range results.List {
		for _, n := range p.Names {
			renames[n.Name] = fmt.Sprintf("%s.%s", varName, n.Name)
		}
	}
	dst.Inspect(fn.Body, func(n dst.Node) bool {
		switch i := n.(type) {
		case *dst.ReturnStmt:
			results := i.Results
			i.Results = []dst.Expr{
				// result
				&dst.CompositeLit{
					Type: &dst.Ident{
						Name: resultName,
					},
					Elts: results,
				},
			}
		}
		return true
	})
	for i, r := range results.List {
		if len(r.Names) < 1 {
			r.Names = []*dst.Ident{&dst.Ident{Name: fmt.Sprintf("Field%d", i)}}
		}
	}
	decl := newStruct(resultName, results)
	f.Decls = append(f.Decls, decl)
	updated := map[string]*dst.File{p.file: f}
	for filename, f := range files {
		dst.Inspect(f, func(n dst.Node) bool {
			if n == nil {
				return true
			}
			switch c := n.(type) {
			default:
				if false {
					// TODO: update references
					fmt.Printf("[%s] node: %T - %+v\n", filename, c, c)
				}
			}
			return true
		})
	}
	// TODO other packages
	return updated, nil
}
