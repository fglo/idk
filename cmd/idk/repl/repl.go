package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fglo/idk/pkg/idk/ast"
	"github.com/fglo/idk/pkg/idk/parser"
)

const PROMPT = "$ "

const IDK_LOGO = `╦╔╦╗╦╔═
║ ║║╠╩╗
╩═╩╝╩ ╩`

// TODO: finish/fix REPL
func Start(in io.Reader, out io.Writer, prettyPrint bool) {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }
	scanner := bufio.NewScanner(in)
	// env := object.NewEnvironment()

	fmt.Print(IDK_LOGO)
	fmt.Println("")
	fmt.Println("")

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "q!" {
			break
		}

		p := parser.NewParser(string(line))
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		} else if prettyPrint {
			ast.PrettyPrintProgram(program)
		}

		// evaluated := evaluator.Eval(program, env)
		// if evaluated != nil {
		// 	io.WriteString(out, evaluated.Inspect())
		// 	io.WriteString(out, "\n")
		// }
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
