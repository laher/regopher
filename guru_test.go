package main

import (
	"strings"
	"testing"
)

func TestParseGuruReferrers(t *testing.T) {

	const jsonStream = `{
        "objpos": "/home/am/go/src/github.com/laher/regopher/extract_parameter_object.go:13:6",
        "desc": "type github.com/laher/regopher.guruPos struct{file string; start int; end int}"
}
{
        "package": "github.com/laher/regopher",
        "refs": [
                {
                        "pos": "/home/am/go/src/github.com/laher/regopher/extract_parameter_object.go:19:48",
                        "text": "func doExtractParameterObject(string name, ref guruPos) {"
                }
        ]
}`
	init, packages, err := parseGuruReferrers(strings.NewReader(jsonStream))
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}
	if init.ObjPos == "" {
		t.Error("expected an objpos for referred object")
	}
	if len(packages) != 1 {
		t.Error("Expected one package")
	}
	p := packages[0]
	if p.Package != "github.com/laher/regopher" {
		t.Error("Expected package to be laher/regopher")
	}
	if len(p.Refs) != 1 {
		t.Error("Expected one reference")
	}
	r := p.Refs[0]

	pos, err := parsePositionString(r.Pos)
	if err != nil {
		t.Errorf("Parse guruPos error: %v", err)
	}
	if pos.file != "/home/am/go/src/github.com/laher/regopher/extract_parameter_object.go" ||
		pos.line != 19 ||
		pos.col != 48 {
		t.Error("Position string parsed incorrectly")
	}
}
