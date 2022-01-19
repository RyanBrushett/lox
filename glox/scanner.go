package glox

import "errors"

type scanner struct {
	source    string
	tokenList []*Token
	start     int
	current   int
	line      int
}

func NewScanner(source string) *scanner {
	var tokenList []*Token

	return &scanner{
		source:    source,
		tokenList: tokenList,
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
			// A comment reaches the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
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
		return s.aString()
	default:
		return RuntimeError(s.line, errors.New("unexpected character"))
	}
}

func (s *scanner) aString() error {
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

func (s *scanner) advance() byte {
	character := s.source[s.current]
	s.current += 1
	return character
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

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) addToken(tt TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	token := NewToken(tt, text, literal, s.line)
	s.tokenList = append(s.tokenList, token)
}
