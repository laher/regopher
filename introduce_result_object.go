package main

import (
	"fmt"

	"github.com/dave/dst"
)

const cmdIntroduceResultObject = "introduce-result-object"

// consolidate results into a struct
// exclude the last value if it's an error
func introduceResultObject(p inputPos, files map[string]*dst.File, fn *dst.FuncDecl) (map[string]*dst.File, error) {
	resultName := fn.Name.Name + "Result"
	varName := "res"
	results := dst.Clone(fn.Type.Results).(*dst.FieldList)
	f := files[p.file]
	var errResult *dst.Ident
	if len(results.List) > 0 {
		if ident, ok := results.List[len(results.List)-1].Type.(*dst.Ident); ok {
			if ident.Name == "error" {
				errResult = ident
				results.List = results.List[:len(results.List)-1]
			}
		}
	}
	fn.Type.Results.List = []*dst.Field{
		&dst.Field{Type: &dst.Ident{Name: resultName}},
	}
	if errResult != nil {
		fn.Type.Results.List = append(fn.Type.Results.List, &dst.Field{Type: errResult})
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
			var errResultVal dst.Expr
			if errResult != nil {
				errResultVal = results[len(results)-1]
				results = results[:len(results)-1]
			}
			i.Results = []dst.Expr{
				// result
				&dst.CompositeLit{
					Type: &dst.Ident{
						Name: resultName,
					},
					Elts: results,
				},
			}
			if errResultVal != nil {
				i.Results = append(i.Results, errResultVal)
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
