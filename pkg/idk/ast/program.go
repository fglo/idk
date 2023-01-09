package ast

import (
	"bytes"
	"fmt"

	"github.com/fglo/idk/pkg/idk/token"
)

type Program struct {
	Statements []Statement
}

func (p *Program) GetTokenValue() string         { return "" }
func (p *Program) GetTokenType() token.TokenType { return "" }
func (p *Program) GetLineNumber() int            { return 0 }
func (p *Program) GetPositionInLine() int        { return 0 }
func (p *Program) GetChildren() []Node           { return []Node{} }
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func PrettyPrintProgram(program *Program) {
	for i, s := range program.Statements {
		fmt.Println(s)
		PrettyPrint(s, "", i == len(program.Statements)-1)
	}
}

func PrettyPrint(node Node, indent string, isLast bool) {

	marker := "├──"
	if isLast {
		marker = "└──"
	}

	fmt.Print(indent)
	fmt.Print(marker)

	fmt.Print(node.GetTokenType())
	fmt.Print(" ")
	fmt.Println(node.GetTokenValue())

	if !isLast {
		indent += "│   "
	} else {
		indent += "    "
	}

	children := node.GetChildren()
	for i, child := range children {
		PrettyPrint(child, indent, i == len(children)-1)
	}
}
