package main

import (
	"flag"
	"os"

	"github.com/fglo/idk/cmd/idk/interpreter"
	"github.com/fglo/idk/cmd/idk/repl"
)

func main() {
	var sourceCodePath string
	var moduleEntryPointPath string
	var prettyPrint bool

	flag.StringVar(&sourceCodePath, "f", "", "File path to the source code.")
	flag.StringVar(&moduleEntryPointPath, "m", "", "File path to the module entry point.")
	flag.BoolVar(&prettyPrint, "p", false, "Pretty print the AST.")
	flag.Parse()

	switch {
	case sourceCodePath != "":
		interpreter.RunSingleFile(sourceCodePath, prettyPrint)
	case moduleEntryPointPath != "":
		interpreter.RunModule(moduleEntryPointPath, prettyPrint)
	default:
		repl.Start(os.Stdin, os.Stdout, prettyPrint)
	}
}
