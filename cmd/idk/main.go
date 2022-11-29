package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fglo/idk/cmd/idk/repl"
	"github.com/fglo/idk/pkg/idk/ast"
	"github.com/fglo/idk/pkg/idk/evaluator"
	"github.com/fglo/idk/pkg/idk/parser"
	"github.com/fglo/idk/pkg/idk/symbol"
)

func run(sourceCodePath string, prettyPrint bool) {
	fileContent, err := os.ReadFile(sourceCodePath)
	check(err)

	p := parser.NewParser(string(fileContent))
	program := p.ParseProgram()
	_ = program

	if len(p.Errors()) != 0 {
		fmt.Println("Parser errors:")
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
	} else {
		if prettyPrint {
			ast.PrettyPrintProgram(program)
		}

		scope := symbol.NewScope()
		result := evaluator.Eval(program, scope)
		if symbol.IsError(result) {
			fmt.Println(result.Inspect())
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var sourceCodePath string
	var prettyPrint bool
	flag.StringVar(&sourceCodePath, "f", "", "Source code file path.")
	flag.BoolVar(&prettyPrint, "p", false, "Pretty print the AST.")
	flag.Parse()

	if sourceCodePath != "" {
		run(sourceCodePath, prettyPrint)
	} else {
		repl.Start(os.Stdin, os.Stdout, prettyPrint)
	}
}
