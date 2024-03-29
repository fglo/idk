package parser

import (
	"fmt"
	"testing"

	"github.com/fglo/idk/pkg/idk/ast"
)

// TODO: more unit tests

func TestDeclareAssignStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x := 1", "x", 1},
		{"y := true", "y", true},
		{"y := false", "y", false},
		{"z := y", "z", "y"},
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testDeclareAssignStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.DeclareAssignStatement).Expression
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"t := -15", "-", 15},
		{"t := -foobar", "-", "foobar"},
		{"t := !foobar", "!", "foobar"},
		{"t := !true", "!", true},
		{"t := !false", "!", false},
	}

	for _, tt := range prefixTests {
		p := NewParser(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.DeclareAssignStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.DeclareAssignStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.GetTokenValue() != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.GetTokenValue())
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"t := 5 + 5", 5, "+", 5},
		{"t := 5 - 5", 5, "-", 5},
		{"t := 5 * 5", 5, "*", 5},
		{"t := 5 / 5", 5, "/", 5},
		{"t := 5 > 5", 5, ">", 5},
		{"t := 5 < 5", 5, "<", 5},
		{"t := 5 == 5", 5, "==", 5},
		{"t := 5 != 5", 5, "!=", 5},
		{"t := foobar + barfoo", "foobar", "+", "barfoo"},
		{"t := foobar - barfoo", "foobar", "-", "barfoo"},
		{"t := foobar * barfoo", "foobar", "*", "barfoo"},
		{"t := foobar / barfoo", "foobar", "/", "barfoo"},
		{"t := foobar > barfoo", "foobar", ">", "barfoo"},
		{"t := foobar < barfoo", "foobar", "<", "barfoo"},
		{"t := foobar == barfoo", "foobar", "==", "barfoo"},
		{"t := foobar != barfoo", "foobar", "!=", "barfoo"},
		{"t := true == true", true, "==", true},
		{"t := true != false", true, "!=", false},
		{"t := false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		p := NewParser(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.DeclareAssignStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.DeclareAssignStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"t := -a * b",
			"((-a) * b)",
		},
		{
			"t := -(a * b)",
			"(-(a * b))",
		},
		{
			"t := !-a",
			"(!(-a))",
		},
		{
			"t := a + b + c",
			"((a + b) + c)",
		},
		{
			"t := a + b - c",
			"((a + b) - c)",
		},
		{
			"t := a * b * c",
			"((a * b) * c)",
		},
		{
			"t := a * b / c",
			"((a * b) / c)",
		},
		{
			"t := a + b / c",
			"(a + (b / c))",
		},
		{
			"t := a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"t := 3 + 4 * -5 * 5",
			"(3 + ((4 * (-5)) * 5))",
		},
		{
			"t := 5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"t := 5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"t := 3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"t := true",
			"true",
		},
		{
			"t := false",
			"false",
		},
		{
			"t := 3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"t := 3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"t := 1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"t := (5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"t := 2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"t := 2 / (5 + 5) % 2",
			"((2 / (5 + 5)) % 2)",
		},
		{
			"t := 2 * 2 % 2 + 2",
			"(((2 * 2) % 2) + 2)",
		},
		{
			"t := (5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"t := -(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"t := !(true == true)",
			"(!(true == true))",
		},
		{
			"t := a + add(b * c, e + f) + d",
			"((a + add((b * c), (e + f))) + d)",
		},
		{
			"t := add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"t := add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		// {
		// 	"t := a * [1, 2, 3, 4][b * c] * d",
		// 	"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		// },
		// {
		// 	"t := add(a * b[2], b[1], 2 * [1, 2][1])",
		// 	"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		// },
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.DeclareAssignStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.DeclareAssignStatement. got=%T",
				program.Statements[0])
		}

		actual := stmt.Expression.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testDeclareAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	declareAssign, ok := s.(*ast.DeclareAssignStatement)
	if !ok {
		t.Errorf("s not *ast.DeclareAssignStatement. got=%T", s)
		return false
	}

	if declareAssign.Identifier.GetValue() != name {
		t.Errorf("declAssStmt.Identifier.Value not '%s'. got=%s", name, declareAssign.Identifier.GetValue())
		return false
	}

	if declareAssign.Identifier.GetValue() != name {
		t.Errorf("declAssStmt.Identifier.GetTokenValue() not '%s'. got=%s", name, declareAssign.Identifier.GetValue())
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.GetValue() != value {
		t.Errorf("integ.Value not %d. got=%d", value, integer.GetValue())
		return false
	}

	if integer.GetTokenValue() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.GetTokenValue() not '%d'. got=%s", value, integer.GetTokenValue())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	identifier, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if identifier.GetValue() != value {
		t.Errorf("ident.Value not %s. got=%s", value, identifier.GetValue())
		return false
	}

	if identifier.GetTokenValue() != value {
		t.Errorf("integer.GetTokenValue() not '%s'. got=%s", value, identifier.GetTokenValue())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if boolean.GetValue() != value {
		t.Errorf("bo.Value not %t. got=%t", value, boolean.GetValue())
		return false
	}

	if boolean.GetTokenValue() != fmt.Sprintf("%t", value) {
		t.Errorf("integer.GetTokenValue() not '%t'. got=%s", value, boolean.GetTokenValue())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.GetTokenValue() != operator {
		t.Errorf("exp.GetTokenValue() is not '%s'. got=%q", operator, opExp.GetTokenValue())
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
