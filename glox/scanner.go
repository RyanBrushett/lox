package glox

import (
	"errors"
	"strconv"
	"unicode"
)

type scanner struct {
	source    string
	tokenList []*Token
	keywords  map[string]TokenType
	start     int
	current   int
	line      int
}

func NewScanner(source string) *scanner {
	var tokenList []*Token
	var keywords = map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}

	return &scanner{
		source:    source,
		tokenList: tokenList,
		keywords:  keywords,
	}
}

func (s *scanner) ScanTokens() ([]*Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}

	s.tokenList = append(
		s.tokenList,
		NewToken(EOF, "", nil, s.line),
	)

	return s.tokenList, nil
}

func (s *scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
		return nil
	case ')':
		s.addToken(RIGHT_PAREN, nil)
		return nil
	case '{':
		s.addToken(LEFT_BRACE, nil)
		return nil
	case '}':
		s.addToken(RIGHT_BRACE, nil)
		return nil
	case ',':
		s.addToken(COMMA, nil)
		return nil
	case '.':
		s.addToken(DOT, nil)
		return nil
	case '-':
		s.addToken(MINUS, nil)
		return nil
	case '+':
		s.addToken(PLUS, nil)
		return nil
	case ';':
		s.addToken(SEMICOLON, nil)
		return nil
	case '*':
		s.addToken(STAR, nil)
		return nil
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
		return nil
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
		return nil
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
		return nil
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
		return nil
	case '/':
		if s.match('/') {
			// A slash-style comment reaches the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.cStyleComment()
		} else {
			s.addToken(SLASH, nil)
		}
		return nil
	case ' ', '\r', '\t':
		return nil
	case '\n':
		s.line += 1
		return nil
	case '"':
		return s.stringLiteral()
	default:
		if isDigit(c) {
			return s.numberLiteral()
		} else if isAlpha(c) {
			s.identifier()
			return nil
		} else {
			return RuntimeError(s.line, errors.New("unexpected character"))
		}
	}
}

func (s *scanner) stringLiteral() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		return RuntimeError(s.line, errors.New("unterminated string"))
	}

	s.advance() // This is the closing '"'
	stringLiteral := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, stringLiteral)
	return nil
}

func (s *scanner) numberLiteral() error {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	if numLiteral, err := strconv.ParseFloat(s.source[s.start:s.current], 64); err == nil {
		s.addToken(NUMBER, numLiteral)
		return nil
	} else {
		return err
	}
}

func (s *scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	foundType, ok := s.keywords[text]
	if ok {
		s.addToken(foundType, nil)
		return
	}

	s.addToken(IDENTIFIER, nil)
}

func (s *scanner) cStyleComment() {
	for !s.isAtEnd() {
		c := s.advance()
		switch c {
		case '\n':
			s.line += 1
		case '/':
			if s.match('*') {
				s.cStyleComment()
			}
		case '*':
			if s.peek() == '/' {
				s.advance()
				return
			}
		default:
			continue
		}
	}
}

func (s *scanner) advance() byte {
	character := s.source[s.current]
	s.current += 1
	return character
}

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current += 1
	return true
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) addToken(tt TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	token := NewToken(tt, text, literal, s.line)
	s.tokenList = append(s.tokenList, token)
}

func isDigit(c byte) bool {
	return unicode.IsNumber(rune(c))
}

func isAlpha(c byte) bool {
	return unicode.IsLetter(rune(c)) || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
