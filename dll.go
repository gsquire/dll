package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"runtime"
	"sync"
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
		os.Exit(1)
	}

	files := os.Args[1:]
	fileCount := len(files)

	reportsChannel := make(chan []*report, fileCount)
	cpuCount := runtime.NumCPU()
	filesPerCore := splitFilesIntoParts(files, cpuCount)

	go func() {
		var wg sync.WaitGroup

		for _, files := range filesPerCore {
			wg.Add(1)

			go func(files []string) {
				defer wg.Done()
				for _, source := range files {
					r, err := gather(source, true)
					if err != nil {
						fmt.Fprintf(os.Stderr, "error parsing %s: %s\n", source, err)
						reportsChannel <- []*report{}
						return
					}
					reportsChannel <- r
				}
			}(files)
		}

		wg.Wait()
		close(reportsChannel)
	}()

	for reports := range reportsChannel {
		for _, report := range reports {
			fmt.Println(report)
		}
	}
}

func splitFilesIntoParts(files []string, parts int) [][]string {
	fileCount := len(files)

	if parts == 1 || fileCount == 1 {
		return [][]string{files}
	}

	for (fileCount % parts) != 0 {
		parts--
	}

	filesPerCPU := fileCount / parts

	fileParts := make([][]string, 0, parts)
	x := 0
	for i := 0; i < parts; i++ {
		fileParts = append(fileParts, files[x:(x+filesPerCPU)])
		x = x + filesPerCPU
	}

	return fileParts
}
