package glox

import "fmt"

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

func (tt TokenType) String() string {
	switch tt {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EOF:
		return "EOF"
	case STRING:
		return "STRING"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	default:
		return "Unknown"
	}
}

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
	return fmt.Sprintf("%v %v %v %v", t.tokenType, t.lexeme, t.literal, t.line)
}
