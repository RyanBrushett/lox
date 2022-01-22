package glox

import "testing"

func TestTokenString(t *testing.T) {
	token := NewToken(EQUAL_EQUAL, "==", nil, 0)
	expected := "EQUAL_EQUAL == <nil> 0"
	actual := token.String()

	if actual != expected {
		t.Errorf("Expected: %s Actual: %s", expected, actual)
	}
}

func TestTokenTypeString(t *testing.T) {
	tt := TokenType(DOT)
	expected := "DOT"
	actual := tt.String()

	if actual != expected {
		t.Errorf("Expected: %s Actual: %s", expected, actual)
	}
}
