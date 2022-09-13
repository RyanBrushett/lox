package glox

import "testing"

func TestParserParsesCorrectlySimple(t *testing.T) {
	parser := simpleTestParser("1 + 1", t)
	expression := parser.parseExpression()
	expected := "(+ 1.0 1.0)"
	actual, _ := NewAstPrinter().print(expression)
	if actual != expected {
		t.Errorf(
			"Incorrect parsing.\n\nExpected: %s\nGot: %s",
			expected, actual,
		)
	}
}

func TestParserParsesCorrectlyAcceptanceTestExample(t *testing.T) {
	parser := simpleTestParser("(5 - (3 - 1)) + -1", t)
	expression := parser.parseExpression()
	expected := "(+ (group (- 5.0 (group (- 3.0 1.0)))) (- 1.0))"
	actual, _ := NewAstPrinter().print(expression)
	if actual != expected {
		t.Errorf(
			"Incorrect parsing.\n\nExpected: %s\nGot: %s",
			expected, actual,
		)
	}
}

func TestParserParsesCorrectlyBetterExample(t *testing.T) {
	parser := simpleTestParser("6 + 3 * 2 / 3 - 1 + -1 + (3 + 3)", t)
	expression := parser.parseExpression()
	expected := "(+ (+ (- (+ 6.0 (* 3.0 (/ 2.0 3.0))) 1.0) (- 1.0)) (group (+ 3.0 3.0)))"
	actual, _ := NewAstPrinter().print(expression)
	if actual != expected {
		t.Errorf(
			"Incorrect parsing.\n\nExpected: %s\nGot: %s",
			expected, actual,
		)
	}
}

func TestParserParsesCorrectlyGreaterLess(t *testing.T) {
	parser := simpleTestParser("6 < 3 <= 3 >= 1 > 0", t)
	expression := parser.parseExpression()
	expected := "(> (>= (<= (< 6.0 3.0) 3.0) 1.0) 0.0)"
	actual, _ := NewAstPrinter().print(expression)
	if actual != expected {
		t.Errorf(
			"Incorrect parsing.\n\nExpected: %s\nGot: %s",
			expected, actual,
		)
	}
}

func TestParserParsesCorrectlyBangEqualNilTrueFalse(t *testing.T) {
	parser := simpleTestParser(`1 != 2 == true == false == nil != "hello"`, t)
	expression := parser.parseExpression()
	expected := "(!= (== (== (== (!= 1.0 2.0) true) false) nil) hello)"
	actual, _ := NewAstPrinter().print(expression)
	if actual != expected {
		t.Errorf(
			"Incorrect parsing.\n\nExpected: %s\nGot: %s",
			expected, actual,
		)
	}
}

func TestParserTernaryOperator(t *testing.T) {
	parser := simpleTestParser(`1 > 2 ? 1 : 2`, t)
	expression := parser.parseExpression()
	expected := `(?: (> 1.0 2.0) 1.0 2.0)`
	actual, _ := NewAstPrinter().print(expression)
	if actual != expected {
		t.Errorf(
			"Incorrect parsing.\n\nExpected: %s\nGot: %s",
			expected, actual,
		)
	}
}

func TestParserTernaryOperatorNested(t *testing.T) {
	parser := simpleTestParser(`1 == 2 ? 3 ? 4 : 5 : 6 ? 7 : 8`, t)
	expression := parser.parseExpression()
	expected := `(?: (== 1.0 2.0) (?: 3.0 4.0 5.0) (?: 6.0 7.0 8.0))`
	actual, _ := NewAstPrinter().print(expression)
	if actual != expected {
		t.Errorf(
			"Incorrect parsing.\n\nExpected: %s\nGot: %s",
			expected, actual,
		)
	}
}

func TestParserNotIsAtEnd(t *testing.T) {
	parser := simpleTestParser("123", t)
	if parser.isAtEnd() {
		t.Errorf("Parser should not be at EOF")
	}
}

func TestParserIsAtEnd(t *testing.T) {
	parser := simpleTestParser("", t)
	if !parser.isAtEnd() {
		t.Errorf("Parser should be at EOF")
	}
}

func TestParserPeek(t *testing.T) {
	parser := simpleTestParser("123", t)

	expected := NewToken(NUMBER, "123", 123.0, 0)
	if parser.peek().String() != expected.String() {
		t.Errorf(
			"Parser did not peek the correct token.\n\nExp: %s\nGot: %s",
			expected.String(), parser.peek().String(),
		)
	}
}

func TestParserPrevious(t *testing.T) {
	parser := simpleTestParser("0 1.0 2.0 3.0", t)
	parser.current += 2
	expected := 1.0
	actual := parser.previous().literal
	if actual != expected {
		t.Errorf(
			"Parser previous token did not return expected value\n\nExp: %v\nGot: %v",
			expected, actual,
		)
	}
}

func TestParserPreviousEmpty(t *testing.T) {
	parser := simpleTestParser("", t)
	actual := parser.previous()
	if actual != nil {
		t.Errorf("Expected previous to be nil but got %v", actual)
	}
}

func TestParserAdvance(t *testing.T) {
	parser := simpleTestParser("0 1.0 2.0 3.0", t)
	expectedPrevious := 0.0
	previous := parser.advance().literal
	expectedPeek := 1.0
	peek := parser.peek().literal

	if previous != expectedPrevious {
		t.Errorf(
			"Previous not what we expected.\n\nExp: %v\nGot: %v",
			expectedPrevious, previous,
		)
	}

	if peek != expectedPeek {
		t.Errorf(
			"Peek not what we expected.\n\nExp: %v\nGot: %v",
			expectedPeek, peek,
		)
	}
}

func TestParserAdvanceAtEnd(t *testing.T) {
	parser := simpleTestParser("", t)
	actual := parser.advance()
	if actual != nil {
		t.Errorf("Expected nil advance for empty source")
	}
}

func TestParserCheckTrue(t *testing.T) {
	parser := simpleTestParser(`"hello" "hello"`, t)
	if !parser.check(STRING) {
		t.Errorf(
			"Expected check of STRING to STRING to return true.",
		)
	}
}

func TestParserCheckFalse(t *testing.T) {
	parser := simpleTestParser("1.0", t)
	if parser.check(STRING) {
		t.Errorf(
			"Expected check of STRING to NUMBER to return false.",
		)
	}
}

func TestParserMatchTrue(t *testing.T) {
	parser := simpleTestParser("<", t)
	match := parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL)

	if !match {
		t.Errorf(`Expected LESS to match from source "<"`)
	}
}

func TestParserMatchFalse(t *testing.T) {
	parser := simpleTestParser("<", t)
	match := parser.match(BANG, BANG_EQUAL)

	if match {
		t.Errorf(`Expected no match to BANG and BANG_EQUAL from source "<"`)
	}
}

func TestParserConsume(t *testing.T) {
	parser := simpleTestParser("()", t)
	_ = parser.advance() // parse the "("
	_, err := parser.consume(RIGHT_PAREN, "Expect ')' after expression")

	if err != nil {
		t.Errorf("Expected consume to return OK")
	}
}

func TestParserConsumeFails(t *testing.T) {
	parser := simpleTestParser("(", t)
	_ = parser.advance() // parse the "("
	_, err := parser.consume(RIGHT_PAREN, "Expect ')' after expression")

	if err == nil {
		t.Errorf("Expected consume to return an error")
	}

	expected := "[Line 0] Error at end: Expect ')' after expression\n"
	actual := err.Error()
	if err.Error() != expected {
		t.Errorf(
			"Error message incorrect.\n\nExp: %s\nGot: %v",
			expected, actual,
		)
	}
}

func simpleTestParser(source string, t *testing.T) *parser {
	tokens, err := NewScanner(source).ScanTokens()

	if err != nil {
		t.Errorf("Tokens did not scan correctly, that's pretty bad")
	}

	return NewParser(tokens)
}
