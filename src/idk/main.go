package main

import (
	"fmt"
	"idk/ast"
	"idk/lexer"
	"idk/parser"
	"idk/repl"
	"os"
)

func run(sourceCodePath string) {
	fileContent, err := os.ReadFile(sourceCodePath)
	check(err)

	l := lexer.NewLexer(string(fileContent))
	_ = l

	// for {
	// 	fmt.Println(l.ReadToken())
	// 	if l.PeekNext() == '\000' {
	// 		break
	// 	}
	// }

	p := parser.NewParser(string(fileContent))
	program := p.ParseProgram()
	_ = program
	ast.PrettyPrintProgram(program)

	if len(p.Errors()) != 0 {
		fmt.Println("Parser errors:")
		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) >= 2 {
		sourceCodePath := os.Args[1]
		run(sourceCodePath)
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
