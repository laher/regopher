package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// flags
var (
	// TODO: modifiedFlag = flag.Bool("modified", false, "read archive of modified files from standard input")
	// TODO: actually do something with scope
	scopeFlag = flag.String("scope", "", "comma-separated list of `packages` the analysis should be limited to")
	writeFlag = flag.Bool("write", false, "write directly to files in place (instead of stdout)")
)

const useHelp = "Run 'regopher -help' for more information.\n"

const helpMessage = `Go source code regophering (refactoring).
Usage: regopher [flags] <mode> <position>

The mode argument determines the query to perform:

	introduce-parameter-object	replace parameters with a struct

The position argument specifies the filename and byte offset (or range)
of the syntax element to query.  For example:

	foo.go:#123,#128
	bar.go:#123

The -write flag causes regopher write back to files in-place; This will be incompatible with -modified
	NOTE: -write might become the default behaviour (except for -modified)

TODO: The -modified flag causes regopher to read an archive from standard input.
	Files in this archive will be used in preference to those in
	the file system.  In this way, a text editor may supply regopher
	with the contents of its unsaved buffers.  Each archive entry
	consists of the file name, a newline, the decimal file size,
	another newline, and the contents of the file.

TODO: The -scope flag restricts analysis to the specified packages.
	Its value is a comma-separated list of patterns of these forms:
		github.com/laher/regopher     # a single package
		github.com/laher/...          # all packages beneath dir
		...                           # the entire workspace.
	A pattern preceded by '-' is negative, so the scope
		encoding/...,-encoding/xml
	matches all encoding packages except encoding/xml.

Example: introduce-parameter-object at offset 530 in this file (an import spec):

  $ regopher introduce-parameter-object main.go:#530
`

/*
inputPos

The position argument specifies the filename and byte offset (or range)
of the syntax element to query.  For example:

	foo.go:#123,#128
	bar.go:#123
*/
type inputPos struct {
	file string
	pos  int
	end  int
}

func parseLineNum(s string) (int, error) {
	if !strings.HasPrefix(s, "#") {
		return 0, fmt.Errorf("missing hash")
	}
	pos, err := strconv.Atoi(s[1:])
	if err != nil {
		return 0, err
	}
	return pos, nil
}
func parseInputPositionString(str string) (inputPos, error) {
	parts := strings.Split(str, ":")
	if len(parts) != 2 {
		return inputPos{}, fmt.Errorf("There should be exactly one : symbol")
	}
	file := parts[0]
	if file != "" {
		file = filepath.Clean(file)
	}
	p := inputPos{file: file}
	posParts := strings.Split(parts[1], ",")
	var err error
	p.pos, err = parseLineNum(posParts[0])
	if err != nil {
		return p, err
	}
	if len(posParts) == 2 {
		p.end, err = parseLineNum(posParts[1])
		if err != nil {
			return p, err
		}
	}
	if len(posParts) > 2 {
		return p, fmt.Errorf("There should be one or 2 line numbers")
	}
	return p, nil
}

type query struct {
	Pos   string         // query position
	Build *build.Context // package loading configuration

	// pointer analysis options
	Scope []string // main packages in (*loader.Config).FromArgs syntax

	// result-printing function, safe for concurrent use
	//Output func(*token.FileSet, QueryResult)
}

func printHelp() {
	fmt.Fprintln(os.Stderr, helpMessage)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
}

func main() {
	log.SetPrefix("regopher: ")
	log.SetFlags(0)

	// Don't print full help unless -help was requested.
	// Just gently remind users that it's there.
	flag.Usage = func() { fmt.Fprint(os.Stderr, useHelp) }
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError) // hack
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		// (err has already been printed)
		if err == flag.ErrHelp {
			printHelp()
		}
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}
	mode, posn := args[0], args[1]

	if mode == "help" {
		printHelp()
		os.Exit(2)
	}

	// Avoid corner case of split("").
	var scope []string
	if *scopeFlag != "" {
		scope = strings.Split(*scopeFlag, ",")
	}
	ctxt := &build.Default
	// Do the refactor.
	q := query{
		Pos:   posn,
		Build: ctxt,
		Scope: scope,
	}
	if err := run(mode, &q); err != nil {
		log.Fatal(err)
	}
}

func loadFiles(p inputPos) (*decorator.Decorator, map[string]*dst.File, error) {
	fset := token.NewFileSet()
	files := map[string]*dst.File{}
	d := decorator.New(fset)
	af, err := parser.ParseFile(fset, p.file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil && af == nil {
		return d, files, err
	}
	f := d.DecorateFile(af)
	files[filepath.Clean(p.file)] = f
	dir := filepath.Dir(p.file)
	matches, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return d, files, err
	}
	for _, match := range matches {
		if match != p.file {
			af, err := parser.ParseFile(fset, match, nil, parser.AllErrors|parser.ParseComments)
			if err != nil && af == nil {
				return d, files, err
			}
			f := d.DecorateFile(af)
			files[filepath.Clean(match)] = f
		}
	}
	return d, files, nil
}

func run(mode string, q *query) error {
	updated := map[string]*dst.File{}
	switch mode {
	case cmdIntroduceParameterObject:
		p, err := parseInputPositionString(q.Pos)
		if err != nil {
			return err
		}
		d, files, err := loadFiles(p)
		if err != nil {
			return err
		}

		funcDecl, err := getFuncAt(d, files[filepath.Clean(p.file)], p.pos)
		if err != nil {
			return err
		}
		updated, err = introduceParameterObject(p, files, funcDecl)
		if err != nil {
			return err
		}
	case cmdIntroduceResultObject:
		p, err := parseInputPositionString(q.Pos)
		if err != nil {
			return err
		}
		d, files, err := loadFiles(p)
		if err != nil {
			return err
		}

		funcDecl, err := getFuncAt(d, files[filepath.Clean(p.file)], p.pos)
		if err != nil {
			return err
		}
		updated, err = introduceResultObject(p, files, funcDecl)
		if err != nil {
			return err
		}
	case "no-op":
		p, err := parseInputPositionString(q.Pos)
		if err != nil {
			return err
		}
		_, updated, err = loadFiles(p)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid subcommand '%s'", mode)
	}

	if *writeFlag {
		for name, f := range updated {
			w, err := os.OpenFile(name, os.O_RDWR|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}

			if err := decorator.Fprint(w, f); err != nil {
				return err
			}
			err = w.Close()
			if err != nil {
				return err
			}
		}
	} else {
		for name, f := range updated {
			// filename
			// length (excluding trailing newline)
			// contents
			// newline
			w := &bytes.Buffer{}
			fmt.Fprintln(os.Stdout, name)
			if err := decorator.Fprint(w, f); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "%d\n", w.Len())
			fmt.Fprintln(os.Stdout, w.String())
		}
	}
	return nil
}
