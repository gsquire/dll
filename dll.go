package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type report struct {
	sourceName string
	lineNumber int
}

func (r *report) String() string {
	return fmt.Sprintf("%s: found defer statement in for loop at line number: %d", r.sourceName, r.lineNumber)
}

type walker struct {
	fset    *token.FileSet
	filePos token.Pos
	reports []*report
}

func newWalker(fset *token.FileSet, filePos token.Pos) *walker {
	return &walker{
		fset:    fset,
		filePos: filePos,
		reports: make([]*report, 0),
	}
}

func (w *walker) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	source := w.fset.File(w.filePos)
	sourceName := source.Name()

	switch ty := n.(type) {
	case *ast.ForStmt:
		for _, stmt := range ty.Body.List {
			if d, ok := stmt.(*ast.DeferStmt); ok {
				w.reports = append(w.reports,
					&report{sourceName: sourceName, lineNumber: source.Line(d.Pos())})
			}
		}
	}

	return w
}

func gather(source string, asFile bool) ([]*report, error) {
	fset := token.NewFileSet()

	var (
		f   *ast.File
		err error
	)

	if asFile {
		f, err = parser.ParseFile(fset, source, nil, 0)
	} else {
		f, err = parser.ParseFile(fset, "", source, 0)
	}
	if err != nil {
		return nil, err
	}

	w := newWalker(fset, f.Pos())
	ast.Walk(w, f)
	return w.reports, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no source files supplied")
	}

	for _, source := range os.Args[1:] {
		r, err := gather(source, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing %s: %s\n", source, err)
			continue
		}
		for _, rep := range r {
			fmt.Println(rep)
		}
	}
}
