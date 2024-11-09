package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strconv"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, parser.Mode(0))
	if err != nil {
		panic(err)
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		length := strconv.Itoa(int(fn.End() - fn.Pos()))
		os.Stdout.Write(append([]byte(length), byte('|')))
		os.Stdout.Write(data[fn.Pos() - 1:fn.End()])
	}
}
