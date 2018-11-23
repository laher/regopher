package main

import (
	"fmt"
	"go/token"
	"os"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func TestLearning1(t *testing.T) {
	code := `package main

func funcy(a int, b int, c int) {
	println(a, b, c)
}

var y = x{}

type x struct {
	a int
}`
	name := "parameterStruct"
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	params := &dst.FieldList{}
	for _, d := range f.Decls {
		//	fmt.Printf("%T: %+v\n", d, d)
		switch x := d.(type) {
		case *dst.FuncDecl:
			fmt.Printf("FUNCDECL - %+v\n", x)
			params = dst.Clone(x.Type.Params).(*dst.FieldList)
			for _, p := range params.List {
				fmt.Printf("PARAM - %+v\n", p)
			}
			x.Type.Params.List = []*dst.Field{
				&dst.Field{
					Names: []*dst.Ident{
						&dst.Ident{Name: name},
					},
					Type: &dst.StructType{},
				},
				// &{Names:[y] Type:<nil> Values:[0xc0000ba140] Decs:{NodeDecs:{Space:None Start:[] End:[]
				//After:None} Assign:[]}}
			}

		}
	}
	/*
		call := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)

		call.Decs.Space = dst.EmptyLine
		call.Decs.After = dst.EmptyLine

		for _, v := range call.Args {
			v := v.(*dst.Ident)
			v.Decs.Space = dst.NewLine
			v.Decs.After = dst.NewLine
		}
	*/
	decl := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: &dst.Ident{Name: name},
				Type: &dst.StructType{
					Fields: params,
				},
			},
		},
		Decs: dst.GenDeclDecorations{
			NodeDecs: dst.NodeDecs{},
		},
	}
	f.Decls = append(f.Decls, decl)

	//		Token: "type",
	if len(os.Args) > 1 && os.Args[1] == "r" {
		for _, d := range f.Decls {
			//fmt.Printf("decls: %T: %+v\n", d, d)
			switch x := d.(type) {
			case *dst.FuncDecl:
				fmt.Printf("FUNCDECL - %+v\n", x)

			case *dst.GenDecl:
				for _, s := range x.Specs {
					switch ts := s.(type) {
					case *dst.ValueSpec:
						fmt.Printf("VALUE - spec values[%T] %+v\n", ts.Values[0], ts.Values[0])
					default:
						fmt.Printf("GENDECL[%T] - spec [%T]: %+v\n", x.Tok, s, s)
					}
				}
			default:
				fmt.Printf("OTHER[%T]: %+v\n", d, d)
			}
		}
	}
	if len(os.Args) > 1 && os.Args[1] == "p" {
		if err := decorator.Print(f); err != nil {
			panic(err)
		}
	}
}
