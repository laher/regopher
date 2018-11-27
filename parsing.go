package main

import (
	"fmt"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func getIdentifierAt(d *decorator.Decorator, f *dst.File, pos int) (*dst.Ident, error) {
	var r *dst.Ident
	dst.Inspect(f, func(n dst.Node) bool {
		if n == nil {
			return false
		}
		an, ok := d.Ast.Nodes[n]
		if !ok {
			// ignores things like FuncType
			return true
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
	if r == nil {
		return nil, fmt.Errorf("ident not found")
	}
	return r, nil
}

func getFuncAt(d *decorator.Decorator, f *dst.File, pos int) (*dst.FuncDecl, error) {
	var r *dst.FuncDecl
	dst.Inspect(f, func(n dst.Node) bool {
		if n == nil {
			return false
		}
		an, ok := d.Ast.Nodes[n]
		if !ok {
			// ignores things like FuncType
			return true
		}
		p := an.Pos()
		e := an.End()
		if pos > int(p) && pos < int(e) {
			fmt.Printf("%T [%d,%d]: %+v\n", n, p, e, n)
			if i, ok := n.(*dst.FuncDecl); ok {
				r = i
			}
			// ok
			return true
		}
		return false
	})
	if r == nil {
		return nil, fmt.Errorf("ident not found")
	}
	return r, nil
}

func getFuncByName(f *dst.File, funcName string) (*dst.FuncDecl, error) {

	for _, d := range f.Decls {
		switch fn := d.(type) {
		case *dst.FuncDecl:
			if fn.Name.Name == funcName {
				return fn, nil
			}
		}
	}
	return nil, fmt.Errorf("func not found")
}
