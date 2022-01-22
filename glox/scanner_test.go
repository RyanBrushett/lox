package glox

import "testing"

func TestNumbers(t *testing.T) {
	source := `123
	123.456
	.456
	123.
	`
	tokenList := scanSource(source, t)
	expectedTokens := []*Token{
		NewToken(NUMBER, "123", 123.0, 0),
		NewToken(NUMBER, "123.456", 123.456, 1),
		NewToken(DOT, ".", nil, 2),
		NewToken(NUMBER, "456", 456.0, 2),
		NewToken(NUMBER, "123", 123.0, 3),
		NewToken(DOT, ".", nil, 3),
		NewToken(EOF, "", nil, 4),
	}

	compareTokens(tokenList, expectedTokens, t)
}

func TestIdentifiers(t *testing.T) {
	source := `andy formless fo _ _123 _abc ab123
	abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_
	`

	tokenList := scanSource(source, t)
	expectedTokens := []*Token{
		NewToken(IDENTIFIER, "andy", nil, 0),
		NewToken(IDENTIFIER, "formless", nil, 0),
		NewToken(IDENTIFIER, "fo", nil, 0),
		NewToken(IDENTIFIER, "_", nil, 0),
		NewToken(IDENTIFIER, "_123", nil, 0),
		NewToken(IDENTIFIER, "_abc", nil, 0),
		NewToken(IDENTIFIER, "ab123", nil, 0),
		NewToken(IDENTIFIER, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_", nil, 1),
		NewToken(EOF, "", nil, 2),
	}

	compareTokens(tokenList, expectedTokens, t)
}

func TestStrings(t *testing.T) {
	source := `""
	"string"
	`

	tokenList := scanSource(source, t)
	expectedTokens := []*Token{
		NewToken(STRING, `""`, ``, 0),
		NewToken(STRING, `"string"`, "string", 1),
		NewToken(EOF, "", nil, 2),
	}

	compareTokens(tokenList, expectedTokens, t)
}

func TestKeywords(t *testing.T) {
	source := `and class else false for fun if nil or return super this true var while`

	tokenList := scanSource(source, t)
	expectedTokens := []*Token{
		NewToken(AND, "and", nil, 0),
		NewToken(CLASS, "class", nil, 0),
		NewToken(ELSE, "else", nil, 0),
		NewToken(FALSE, "false", nil, 0),
		NewToken(FOR, "for", nil, 0),
		NewToken(FUN, "fun", nil, 0),
		NewToken(IF, "if", nil, 0),
		NewToken(NIL, "nil", nil, 0),
		NewToken(OR, "or", nil, 0),
		NewToken(RETURN, "return", nil, 0),
		NewToken(SUPER, "super", nil, 0),
		NewToken(THIS, "this", nil, 0),
		NewToken(TRUE, "true", nil, 0),
		NewToken(VAR, "var", nil, 0),
		NewToken(WHILE, "while", nil, 0),
		NewToken(EOF, "", nil, 0),
	}

	compareTokens(tokenList, expectedTokens, t)
}

func TestPunctuators(t *testing.T) {
	source := `(){};,+-*!===<=>=!=<>/.`

	tokenList := scanSource(source, t)
	expectedTokens := []*Token{
		NewToken(LEFT_PAREN, "(", nil, 0),
		NewToken(RIGHT_PAREN, ")", nil, 0),
		NewToken(LEFT_BRACE, "{", nil, 0),
		NewToken(RIGHT_BRACE, "}", nil, 0),
		NewToken(SEMICOLON, ";", nil, 0),
		NewToken(COMMA, ",", nil, 0),
		NewToken(PLUS, "+", nil, 0),
		NewToken(MINUS, "-", nil, 0),
		NewToken(STAR, "*", nil, 0),
		NewToken(BANG_EQUAL, "!=", nil, 0),
		NewToken(EQUAL_EQUAL, "==", nil, 0),
		NewToken(LESS_EQUAL, "<=", nil, 0),
		NewToken(GREATER_EQUAL, ">=", nil, 0),
		NewToken(BANG_EQUAL, "!=", nil, 0),
		NewToken(LESS, "<", nil, 0),
		NewToken(GREATER, ">", nil, 0),
		NewToken(SLASH, "/", nil, 0),
		NewToken(DOT, ".", nil, 0),
		NewToken(EOF, "", nil, 0),
	}

	compareTokens(tokenList, expectedTokens, t)
}

func TestWhitespace(t *testing.T) {
	source := "space     tabs\t\t\t\t\tnewlines\n\n\n\n\nend"

	tokenList := scanSource(source, t)
	expectedTokens := []*Token{
		NewToken(IDENTIFIER, "space", nil, 0),
		NewToken(IDENTIFIER, "tabs", nil, 0),
		NewToken(IDENTIFIER, "newlines", nil, 0),
		NewToken(IDENTIFIER, "end", nil, 5),
		NewToken(EOF, "", nil, 5),
	}

	compareTokens(tokenList, expectedTokens, t)
}

func scanSource(source string, t *testing.T) []*Token {
	scanner := NewScanner(source)
	tokenList, err := scanner.ScanTokens()
	if err != nil {
		t.Errorf("ScanTokens returned an error")
	}

	return tokenList
}

func compareTokens(tokenList []*Token, expectedTokens []*Token, t *testing.T) {
	if len(tokenList) != len(expectedTokens) {
		t.Errorf("Tokens Not Correct. Expected: %d Got: %d", len(expectedTokens), len(tokenList))
	}

	for i, token := range tokenList {
		if *token != *expectedTokens[i] {
			t.Errorf("Token does not match expected. Expected: %v Got %v", token, expectedTokens[i])
		}
	}
}
