package lexer

import (
	"unicode"

	"github.com/fglo/idk/pkg/idk/token"
)

type Lexer struct {
	input          string
	readPosition   int
	position       int
	current        byte
	currentLine    int
	positionInLine int

	errors []string
}

func NewLexer(txt string) *Lexer {
	l := new(Lexer)
	l.input = txt
	l.readPosition = 0
	l.position = -1
	l.current = 0
	l.currentLine = 1
	l.positionInLine = 0
	return l
}

func (l *Lexer) peek(offset int) byte {
	i := l.readPosition + offset - 1
	if i >= len(l.input) {
		return 0
	}
	return l.input[i]
}

func (l *Lexer) PeekNext() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readChar() byte {
	if l.readPosition >= len(l.input) {
		l.current = 0
	} else {
		l.current = l.input[l.readPosition]
		l.position = l.readPosition
		l.positionInLine++
		l.readPosition++
	}
	return l.current
}

func (l *Lexer) skipWhitespace() {
	for ch := rune(l.PeekNext()); unicode.IsSpace(ch) && ch != '\n'; ch = rune(l.PeekNext()) {
		l.readChar()
	}
}

func (l *Lexer) skipEol() {
	l.skipWhitespace()
	for l.PeekNext() == '\n' {
		l.readChar()
		l.currentLine++
		l.skipWhitespace()
	}
}

func (l *Lexer) ReadToken() token.Token {
	l.skipWhitespace()

	ch := l.readChar()

	var tok *token.Token

	tok = token.NewTokenNotDefaultValue(token.ILLEGAL, l.position, l.currentLine, l.positionInLine, string(rune(ch)))
	switch ch {
	case 0:
		tok = token.NewToken(token.EOF, l.position, l.currentLine, l.positionInLine)
	case '\n':
		endlineLine := l.currentLine
		l.skipEol()
		tok = token.NewToken(token.EOL, l.position, endlineLine, l.positionInLine)
		l.currentLine++
		l.positionInLine = 0
	case '+':
		tok = token.NewToken(token.PLUS, l.position, l.currentLine, l.positionInLine)
	case '-':
		tok = token.NewToken(token.MINUS, l.position, l.currentLine, l.positionInLine)
	case '*':
		tok = token.NewToken(token.ASTERISK, l.position, l.currentLine, l.positionInLine)
	case '/':
		if l.PeekNext() == '/' {
			tok = token.NewToken(token.LINE_COMMENT, l.position, l.currentLine, l.positionInLine)
			l.readChar()
		} else {
			tok = token.NewToken(token.SLASH, l.position, l.currentLine, l.positionInLine)
		}
	case '(':
		tok = token.NewToken(token.LPARENTHESIS, l.position, l.currentLine, l.positionInLine)
	case ')':
		tok = token.NewToken(token.RPARENTHESIS, l.position, l.currentLine, l.positionInLine)
	case ':':
		if l.PeekNext() == '=' {
			tok = token.NewToken(token.DECLARE_ASSIGN, l.position, l.currentLine, l.positionInLine)
			l.readChar()
		} else {
			tok = token.NewToken(token.DECLARE, l.position, l.currentLine, l.positionInLine)
		}
	case '=':
		if l.PeekNext() == '=' {
			tok = token.NewToken(token.EQ, l.position, l.currentLine, l.positionInLine)
			l.readChar()
		} else {
			tok = token.NewToken(token.ASSIGN, l.position, l.currentLine, l.positionInLine)
		}
	case '!':
		if l.PeekNext() == '=' {
			tok = token.NewToken(token.NEQ, l.position, l.currentLine, l.positionInLine)
			l.readChar()
		} else {
			tok = token.NewToken(token.NOT, l.position, l.currentLine, l.positionInLine)
		}
	case '<':
		if l.PeekNext() == '=' {
			tok = token.NewToken(token.LTE, l.position, l.currentLine, l.positionInLine)
			l.readChar()
		} else {
			tok = token.NewToken(token.LT, l.position, l.currentLine, l.positionInLine)
		}
	case '>':
		if l.PeekNext() == '=' {
			tok = token.NewToken(token.GTE, l.position, l.currentLine, l.positionInLine)
			l.readChar()
		} else {
			tok = token.NewToken(token.GT, l.position, l.currentLine, l.positionInLine)
		}
	case '.':
		if l.PeekNext() == '.' {
			tok = token.NewToken(token.RANGE, l.position, l.currentLine, l.positionInLine)
			l.readChar()
		}
	case ',':
		tok = token.NewToken(token.COMMA, l.position, l.currentLine, l.positionInLine)
	case '\'':
		tok = l.readCharToken()
	case '"':
		tok = l.readStringToken()
	default:
		switch {
		case unicode.IsDigit(rune(ch)):
			tok = l.readNumberToken()
		case unicode.IsLetter(rune(ch)) || ch == '_':
			tok = l.readWordToken()
		}
	}

	return *tok
}

func (l *Lexer) readNumberToken() *token.Token {
	start := l.position
	startInLine := l.positionInLine
	for ch := rune(l.PeekNext()); unicode.IsDigit(ch); ch = rune(l.PeekNext()) {
		l.readChar()
	}
	number := substring(l.input, start, l.readPosition)
	return token.NewTokenNotDefaultValue(token.INT, start, l.currentLine, startInLine, number)
}

func (l *Lexer) readWordToken() *token.Token {
	start := l.position
	startInLine := l.positionInLine
	for ch := rune(l.PeekNext()); unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'; ch = rune(l.PeekNext()) {
		l.readChar()
	}

	word := substring(l.input, start, l.readPosition)
	keyword := token.LookupKeyword(word)
	return token.NewTokenNotDefaultValue(keyword, start, l.currentLine, startInLine, word)
}

func (l *Lexer) readCharToken() *token.Token {
	if l.peek(2) != '\'' {
		panic("not really a character") // TODO: proper lexer errors
	}

	start := l.readPosition
	startInLine := l.positionInLine
	for ch := l.peek(1); ch != '\''; ch = l.peek(1) {
		l.readChar()
	}
	char := substring(l.input, start, l.readPosition)
	l.readChar()
	return token.NewTokenNotDefaultValue(token.CHAR, start, l.currentLine, startInLine, char)
}

func (l *Lexer) readStringToken() *token.Token {
	start := l.readPosition
	startInLine := l.positionInLine
	for ch := l.peek(1); ch != '"'; ch = l.peek(1) {
		l.readChar()
	}
	str := substring(l.input, start, l.readPosition)
	l.readChar()
	return token.NewTokenNotDefaultValue(token.STRING, start, l.currentLine, startInLine, str)
}

func substring(s string, start, end int) string {
	if start == end {
		return string(s[start])
	}
	return string(s[start:end])
}
