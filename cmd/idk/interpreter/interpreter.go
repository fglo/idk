package interpreter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fglo/idk/pkg/idk/ast"
	"github.com/fglo/idk/pkg/idk/evaluator"
	"github.com/fglo/idk/pkg/idk/parser"
	"github.com/fglo/idk/pkg/idk/symbol"
)

func RunSingleFile(sourceCodePath string, prettyPrint bool) {
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

func RunModule(moduleEntryPoint string, prettyPrint bool) {
	moduleDir := filepath.Dir(moduleEntryPoint)

	files, err := ioutil.ReadDir(moduleDir)
	if err != nil {
		log.Fatal(err)
	}

	scope := symbol.NewScope()
	for _, packageDir := range files {
		if packageDir.IsDir() {
			packageName := packageDir.Name()
			packageDirPath := fmt.Sprintf("%s/%s", moduleDir, packageName)

			packageFiles, err := ioutil.ReadDir(packageDirPath)
			if err != nil {
				log.Fatal(err)
			}

			for _, fileName := range packageFiles {
				filepath := fmt.Sprintf("%s/%s/%s", moduleDir, packageName, fileName.Name())
				parseAndEvalPackageFile(scope, packageName, filepath, prettyPrint)
			}
		}
	}

	parseAndEvalModuleEntryFile(scope, moduleEntryPoint, prettyPrint)
}

func parseAndEvalPackageFile(moduleScope *symbol.Scope, packageName string, filepath string, prettyPrint bool) {
	fileContent, err := os.ReadFile(filepath)
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

		scope := moduleScope.GetNamedScope(packageName)
		result := evaluator.Eval(program, scope)
		if symbol.IsError(result) {
			fmt.Println(result.Inspect())
		}
	}
}

func parseAndEvalModuleEntryFile(scope *symbol.Scope, moduleEntryPoint string, prettyPrint bool) {
	fileContent, err := os.ReadFile(moduleEntryPoint)
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
