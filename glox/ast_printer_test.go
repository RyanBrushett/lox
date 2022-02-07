package glox

import (
	"testing"
)

func TestAstPrinter(t *testing.T) {
	// -123 * (456 + 789)
	expr := NewBinary(
		NewUnary(
			NewToken(MINUS, "-", nil, 1),
			NewLiteral(123),
		),

		NewToken(STAR, "*", nil, 1),

		NewGrouping(
			NewBinary(
				NewLiteral(456),
				NewToken(PLUS, "+", nil, 1),
				NewLiteral(789),
			),
		),
	)

	astPrinter := NewAstPrinter()
	astPrinterOut, err := astPrinter.print(expr)
	checkPrinterError(err, t)

	expected := "(* (- 123) (group (+ 456 789)))"
	assertEqual(expected, astPrinterOut, t)
}

func TestAstPrinterWithNil(t *testing.T) {
	expr := NewLiteral(nil)
	astPrinter := NewAstPrinter()
	astPrinterOut, err := astPrinter.print(expr)
	checkPrinterError(err, t)

	expected := "nil"
	assertEqual(expected, astPrinterOut, t)
}

func checkPrinterError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("astPrinter error: %v", err)
	}
}

func assertEqual(expected string, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Errorf("astPrinter incorrect output:\n\n Expected: %s\nGot: %s", expected, actual)
	}
}
