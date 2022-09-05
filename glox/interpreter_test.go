package glox

import "testing"

func TestInterpreterEvaluatesWholeExpressions(t *testing.T) {
	interpreter := NewInterpreter()
	testcases := map[string]float64{
		"(5 - (3 - 1)) + -1":      2.0,
		"-(1 + 1)":                -2.0,
		"1 == -(2 - 3) ? 1.0 : 0": 1.0,
	}

	for source, expected := range testcases {
		e := parsedExpression(source)
		result, _ := interpreter.evaluate(e)
		assertEqualWithError(result, expected, t, source)
	}
}

func TestInterpreterEvaluatesTruthyBinary(t *testing.T) {
	interpreter := NewInterpreter()
	sources := []string{
		"3 >= 3",
		"2 > 1",
		"1 < 2",
		"3 <= 4",
		"1 == 1",
		"2 != 1",
	}

	for _, source := range sources {
		e := parsedExpression(source)
		result, _ := interpreter.visitBinaryExpr(e.(*Binary))
		if result != true {
			t.Errorf("Expression '%s' should evaluate to true", source)
		}
	}
}

func TestIntegerMath(t *testing.T) {
	interpreter := NewInterpreter()
	testcases := map[string]float64{
		"3 + 3": 6.0,
		"5 - 3": 2.0,
		"6 / 2": 3.0,
		"2 * 3": 6.0,
	}

	for source, expected := range testcases {
		e := parsedExpression(source)
		result, _ := interpreter.visitBinaryExpr(e.(*Binary))
		assertEqualWithError(result, expected, t, source)
	}
}

func TestUnaryExpressions(t *testing.T) {
	interpreter := NewInterpreter()
	testcases := map[string]interface{}{
		"!true": false,
		"-3":    -3.0,
	}

	for source, expected := range testcases {
		e := parsedExpression(source)
		result, _ := interpreter.visitUnaryExpr(e.(*Unary))
		assertEqualWithError(result, expected, t, source)
	}
}

func TestTernaryExpressions(t *testing.T) {
	interpreter := NewInterpreter()
	testcases := map[string]interface{}{
		"2.0 > 1.0 ? 2.0 : 1.0":     2.0,
		"1.0 > 2.0 ? 2.0 : 1.0":     1.0,
		"1.0 == 1.0 ? true : false": true,
		"1.0 == 2.0 ? true : false": false,
	}

	for source, expected := range testcases {
		e := parsedExpression(source)
		result, _ := interpreter.visitTernaryExpr(e.(*Ternary))
		assertEqualWithError(result, expected, t, source)
	}
}

func TestGroupingExpressions(t *testing.T) {
	interpreter := NewInterpreter()
	source := "(3 - 1)"
	e := parsedExpression(source)
	result, _ := interpreter.visitGroupingExpr(e.(*Grouping))
	assertEqualWithError(result, 2.0, t, source)
}

func TestLiteralExpression(t *testing.T) {
	interpreter := NewInterpreter()
	source := "2"
	e := parsedExpression(source)
	result, _ := interpreter.visitLiteralExpr(e.(*Literal))
	assertEqualWithError(result, 2.0, t, source)
}

func parsedExpression(source string) Expr {
	tokens, _ := NewScanner(source).ScanTokens()
	parser := NewParser(tokens)
	return parser.parse()
}

func assertEqualWithError(result interface{}, expected interface{}, t *testing.T, source string) {
	if result != expected {
		t.Errorf(
			"Expression '%s' evaluated incorrectly.\n\nExpected: %v\nGot: %v",
			source, expected, result,
		)
	}
}
