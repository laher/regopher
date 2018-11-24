package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/dave/dst/decorator"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestExtractParameterObject(t *testing.T) {
	testCases := []struct {
		filename string
		function string
	}{
		{
			filename: "parameter_obj_basic.go",
			function: "extractParamUnused",
		},
		{
			filename: "parameter_obj_used.go",
			function: "extractParamUsed",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.filename, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := decorator.ParseFile(fset, "testdata/before/"+testCase.filename, nil, parser.AllErrors)
			if err != nil {
				t.Fatal(err)
			}
			err = extractParameterObject(f, testCase.function)
			if err != nil {
				t.Fatal(err)
			}
			w := &bytes.Buffer{}
			if err := decorator.Fprint(w, f); err != nil {
				t.Fatal(err)
			}
			actual := string(w.Bytes())
			expected, err := ioutil.ReadFile("testdata/expected/" + testCase.filename)
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
