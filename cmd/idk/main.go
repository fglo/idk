package main

import (
	"flag"
	"fmt"
	"github.com/fglo/idk/cmd/idk/repl"
	"github.com/fglo/idk/pkg/idk/ast"
	"github.com/fglo/idk/pkg/idk/evaluator"
	"github.com/fglo/idk/pkg/idk/parser"
	"github.com/fglo/idk/pkg/idk/symbol"
	"os"
)

func run(sourceCodePath string, prettyPrint bool) {
	fileContent, err := os.ReadFile(sourceCodePath)
	check(err)

	// l := lexer.NewLexer(string(fileContent))
	// _ = l

	// for {
	// 	fmt.Println(l.ReadToken())
	// 	if l.PeekNext() == '\000' {
	// 		break
	// 	}
	// }

	p := parser.NewParser(string(fileContent))
	program := p.ParseProgram()
	_ = program

	if len(p.Errors()) != 0 {
		fmt.Println("Parser errors:")
		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg)
		}
	} else if prettyPrint {
		ast.PrettyPrintProgram(program)
	}

	scope := symbol.NewScope()
	evaluated := evaluator.Eval(program, scope)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
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
