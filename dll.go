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
	return fmt.Sprintf("%s:%d found defer statement in for loop", r.sourceName, r.lineNumber)
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
	filesPerCore := splitArrayIntoParts(files, cpuCount)

	go func() {
		var wg sync.WaitGroup

		for _, files := range filesPerCore {
			wg.Add(1)
			go analyseFiles(files, reportsChannel, &wg)
		}

		wg.Wait()
		close(reportsChannel)
	}()

	for reports := range reportsChannel {
		printReports(reports)
	}
}

func analyseFiles(files []string, c chan []*report, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, source := range files {
		r, err := gather(source, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing %s: %s\n", source, err)
			c <- []*report{}
			return
		}
		c <- r
	}
}

func printReports(reports []*report) {
	for _, report := range reports {
		fmt.Println(report)
	}
}

func splitArrayIntoParts(array []string, parts int) [][]string {
	arraySize := len(array)

	if parts <= 1 {
		return [][]string{array}
	}

	// if there are more parts than strings, it tries to find the next smallest number to destribute them equally.
	if arraySize < parts {
		for (arraySize % parts) != 0 {
			parts--
		}
	}

	stringsPerPart := arraySize / parts
	arrayParts := make([][]string, 0, parts)
	lastIndex := 0
	for i := 0; i < parts; i++ {
		arrayParts = append(arrayParts, array[lastIndex:(lastIndex+stringsPerPart)])
		lastIndex = lastIndex + stringsPerPart
	}

	// if not all strings could be splitted equally it will adds the missing ones to the first part.
	if stringsPerPart*parts != arraySize {
		firstpart := arrayParts[0]
		firstpart = append(firstpart, array[stringsPerPart*parts:]...)
		arrayParts[0] = firstpart
	}

	return arrayParts
}
