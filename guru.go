package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/tools/cmd/guru/serial"
)

/*
{
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
}
*/

func runGuruReferrers(pos string) (*serial.ReferrersInitial, []*serial.ReferrersPackage, error) {
	reader, wait, err := invokeGuru("referrers", pos, "-json")
	if err != nil {
		return nil, nil, err
	}
	defer wait()
	ri, rp, err := parseGuruReferrers(reader)
	return ri, rp, err
}

type wait func() error

func invokeGuru(subcommand string, pos string, flags ...string) (io.Reader, wait, error) {
	exe := "guru"
	p, err := exec.LookPath(exe)
	if err != nil {
		log.Printf("Couldn't find exe %s - %s", p, err)
		return nil, nil, err
	}
	cmd := exec.Command(exe)

	cmd.Args = append(cmd.Args, flags...)
	cmd.Args = append(cmd.Args, []string{subcommand, pos}...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Couldn't attach to stdout: %s", err)
		return stdout, nil, err
	}
	err = cmd.Start()
	//log.Printf("args: %v", cmd.Args)
	return stdout, cmd.Wait, err
}

type guruPos struct {
	file string
	line int
	col  int
}

//	file:line:column    valid position with file name
//	file:line           valid position with file name but no column (column == 0)
//	line:column         valid position without file name
//	line                valid position without file name and no column (column == 0)
//	file                invalid position with file name
//	-                   invalid position without file name
//
func parseGuruPos(str string) (guruPos, error) {
	parts := strings.Split(str, ":")
	var (
		p   guruPos
		err error
	)
	switch len(parts) {
	case 1:
		// numeric - line number
		p.line, err = strconv.Atoi(parts[0])
		return p, err
	case 2:
		p.line, err = strconv.Atoi(parts[0])
		if err != nil {
			// str/numeric: file:line
			p.file = parts[0]
			p.line, err = strconv.Atoi(parts[1])
			return p, err
		}
		// numeric/numeric: line:col
		p.col, err = strconv.Atoi(parts[1])
	case 3:
		// str/numeric/numeric: file:line:col
		p.file = parts[0]
		p.line, err = strconv.Atoi(parts[1])
		if err != nil {
			return p, err
		}
		p.col, err = strconv.Atoi(parts[2])
	default:
		err = errors.New("invalid guru Pos string")
	}
	return p, err
}

func parseGuruReferrers(jsonStream io.Reader) (*serial.ReferrersInitial, []*serial.ReferrersPackage, error) {
	dec := json.NewDecoder(jsonStream)
	init := &serial.ReferrersInitial{}
	second := []*serial.ReferrersPackage{}
	err := dec.Decode(init)
	if err != nil {
		return init, second, err
	}
	for dec.More() {
		p := &serial.ReferrersPackage{}
		err := dec.Decode(p)
		if err != nil {
			return init, second, err
		}
		second = append(second, p)
	}
	return init, second, nil
}
