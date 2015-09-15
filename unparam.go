package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/scanner"
)

type UnusedFuncArgsVisitor struct{}

var fset *token.FileSet

func (v *UnusedFuncArgsVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.FuncDecl:
		fmt.Printf("Function name: %q\n", t.Name)
		var params []string
		for _, v := range t.Type.Params.List {
			for _, k := range v.Names {
				params = append(params, k.Name)
			}
		}
		toks := []string{}
		for _, n := range t.Body.List {
			buf := new(bytes.Buffer)
			format.Node(buf, fset, n)
			var s scanner.Scanner
			s.Init(strings.NewReader(buf.String()))
			var tok rune
			for tok != scanner.EOF {
				if tok == scanner.Ident {
					toks = append(toks, s.TokenText())
				}
				tok = s.Scan()
			}
		}
		for _, p := range params {
			var found bool
			for _, v := range toks {
				if v == p {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("WTF, param %q not found\n", p)
			}
		}
	}

	return v
}

func main() {
	fset = token.NewFileSet()
	fileList := []string{}
	err := filepath.Walk(os.Args[1], func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, file := range fileList {
		fmt.Printf("File %s\n", file)
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			continue
		}
		ast.Walk(new(UnusedFuncArgsVisitor), f)
		fmt.Printf("\n")
	}
}
