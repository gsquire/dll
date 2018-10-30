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

type visitFunc func(ast.Node) ast.Visitor

func (vf visitFunc) Visit(node ast.Node) ast.Visitor {
	return vf(node)
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

	var (
		reports             []*report
		findLoop, findDefer ast.Visitor
	)

	findLoop = visitFunc(func(n ast.Node) ast.Visitor {
		if n == nil {
			return nil
		}
		switch n.(type) {
		case *ast.RangeStmt, *ast.ForStmt:
			return findDefer
		default:
			return findLoop
		}
	})

	findDefer = visitFunc(func(n ast.Node) ast.Visitor {
		if n == nil {
			return nil
		}
		switch n := n.(type) {
		case *ast.DeferStmt:
			source := fset.File(f.Pos())
			reports = append(reports, &report{
				sourceName: source.Name(),
				lineNumber: source.Line(n.Pos()),
			})
			return findDefer
		case *ast.FuncLit:
			return findLoop
		default:
			return findDefer
		}
	})

	ast.Walk(findLoop, f)

	return reports, nil
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
