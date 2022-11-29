package token

import (
	"fmt"
	"testing"
)

func TestLookupKeyword(t *testing.T) {
	tests := []struct {
		word string
		want TokenType
	}{
		{"true", BOOL},
		{"false", BOOL},
		{"if", IF},
		{"else", ELSE},
		{"for", FOR},
		{"end", END},
		{"not", NOT},
		{"and", AND},
		{"or", OR},
		{"xor", XOR},
		{"in", IN},
		{"return", RETURN},
		{"func", FUNC},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("testing %s keyword lookup", tt.word), func(t *testing.T) {
			if got := LookupKeyword(tt.word); got != tt.want {
				t.Errorf("LookupKeyword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenType_IsOperator(t *testing.T) {
	tests := []struct {
		tr   TokenType
		want bool
	}{
		{PLUS, true},
		{MINUS, true},
		{ASTERISK, true},
		{SLASH, true},
		{PLUS, true},
		{PLUS, true},
		{GT, true},
		{LT, true},
		{GTE, true},
		{LTE, true},
		{NOT, true},
		{AND, true},
		{OR, true},
		{XOR, true},

		{EOL, false},
		{EOF, false},
		{INT, false},
		{STRING, false},
		{ARRAY, false},
		{BOOL, false},
		{TRUE, false},
		{FALSE, false},
		{DECLASSIGN, false},
		{RANGE, false},
		{RANGE_INCLUSIVE, false},
		{LPARENTHESIS, false},
		{RPARENTHESIS, false},
		{IF, false},
		{ELSE, false},
		{FOR, false},
		{END, false},
		{IN, false},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%s is an operator", tt.tr)
		if !tt.want {
			name = fmt.Sprintf("%s is not an operator", tt.tr)
		}
		t.Run(name, func(t *testing.T) {
			if got := tt.tr.IsOperator(); got != tt.want {
				t.Errorf("TokenType.IsOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}
