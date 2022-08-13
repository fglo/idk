package parser

import (
	"fmt"
	"testing"

	"github.com/fglo/idk/pkg/idk/ast"
)

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

//TODO: lone expressions (need to test them and to have them in the REPL)

// func TestIdentifierExpression(t *testing.T) {
// 	input := "test_name"

// 	p := NewParser(input)
// 	program := p.ParseProgram()
// 	checkParserErrors(t, p)

// 	if len(program.Statements) != 1 {
// 		t.Fatalf("program has not enough statements. got=%d",
// 			len(program.Statements))
// 	}
// 	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
// 	if !ok {
// 		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
// 			program.Statements[0])
// 	}

// 	testIdentifier(t, stmt.Expression, input)
// }

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		p := NewParser(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.UnaryExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Token.Value != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Token.Value)
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
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		p := NewParser(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
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
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		p := NewParser(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
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

	if declareAssign.Identifier.Value != name {
		t.Errorf("declAssStmt.Identifier.Value not '%s'. got=%s", name, declareAssign.Identifier.Value)
		return false
	}

	if declareAssign.Identifier.Token.Value != name {
		t.Errorf("declAssStmt.Identifier.Token.Value not '%s'. got=%s", name, declareAssign.Identifier.Token.Value)
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

	if integer.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.Token.Value != fmt.Sprintf("%d", value) {
		t.Errorf("integer.Token.Value not '%d'. got=%s", value, integer.Token.Value)
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

	if identifier.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, identifier.Value)
		return false
	}

	if identifier.Token.Value != value {
		t.Errorf("integer.Token.Value not '%s'. got=%s", value, identifier.Token.Value)
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

	if boolean.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.Token.Value != fmt.Sprintf("%t", value) {
		t.Errorf("integer.Token.Value not '%t'. got=%s", value, boolean.Token.Value)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.BinaryExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Token.Value != operator {
		t.Errorf("exp.Token.Value is not '%s'. got=%q", operator, opExp.Token.Value)
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
