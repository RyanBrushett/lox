package glox

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	// Single-character tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

//go:generate stringer -type=TokenType

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func NewToken(tt TokenType, lexeme string, lit interface{}, line int) *Token {
	return &Token{
		tokenType: tt,
		lexeme:    lexeme,
		literal:   lit,
		line:      line,
	}
}

func (t *Token) String() string {
	var literal interface{}

	if t.literal == nil {
		literal = "null"
	} else {
		literal = t.literal
	}

	if t.tokenType == NUMBER {
		if strings.Contains(t.lexeme, ".") { // number will already be printed as floating point
			literal = t.literal
		} else {
			literal = fmt.Sprintf("%.1f", t.literal) // lox acceptance tests expect integers to end with `.0`
		}
	}

	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, literal)
}
