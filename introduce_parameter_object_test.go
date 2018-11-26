package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestExtractParameterObject(t *testing.T) {
	testCases := []struct {
		pos string
	}{
		{
			pos: "parameter_obj_basic.go:#42",
		},
		{
			pos: "parameter_obj_used.go:#35",
		},
		{
			pos: "parameter_obj_referenced.go:#35",
		},
		{
			pos: "parameter_obj_method.go:#70",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.pos, func(t *testing.T) {
			fset := token.NewFileSet()
			pos, err := parseInputPositionString(testCase.pos)
			if err != nil {
				t.Fatal(err)
			}
			file := "testdata/before/" + pos.file
			af, err := parser.ParseFile(fset, file, nil, parser.AllErrors|parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			d := decorator.New(fset)
			f, err := d.DecorateFile(af)
			if err != nil {
				t.Fatal(err)
			}
			funcDecl, err := getFuncAt(d, f, pos.pos)
			if err != nil {
				t.Fatal(err)
			}
			_, err = regopherParamsToStruct(pos, map[string]*dst.File{pos.file: f}, funcDecl)
			if err != nil {
				t.Fatal(err)
			}
			w := &bytes.Buffer{}
			if err := decorator.Fprint(w, f); err != nil {
				t.Fatal(err)
			}
			actual := string(w.Bytes())
			expected, err := ioutil.ReadFile("testdata/expected/" + pos.file)
			if err != nil {
				t.Fatal(err)
			}
			if string(expected) != actual {
				dmp := diffmatchpatch.New()

				diffs := dmp.DiffMain(string(expected), actual, false)

				t.Error(dmp.DiffPrettyText(diffs))
			}
		})
	}
}
