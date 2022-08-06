package glox

import (
	"bytes"
	"fmt"
)

type astPrinter struct{}

func NewAstPrinter() *astPrinter {
	return &astPrinter{}
}

func (ap *astPrinter) print(expr Expr) (interface{}, error) {
	return expr.Accept(ap)
}

func (ap *astPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	return ap.parenthesize(expr.Operator.lexeme, expr.Left, expr.Right)
}

func (ap *astPrinter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return ap.parenthesize("group", expr.Expression)
}

func (ap *astPrinter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	switch expr.Value.(type) {
	default:
		return fmt.Sprintf("%v", expr.Value), nil
	case float64:
		if expr.Value == float64(int64(expr.Value.(float64))) {
			return fmt.Sprintf("%.1f", expr.Value), nil // lox acceptance tests expect printed integers to end with `.0`
		} else {
			return fmt.Sprintf("%v", expr.Value), nil
		}
	}
}

func (ap *astPrinter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	return ap.parenthesize(expr.Operator.lexeme, expr.Right)
}

func (ap *astPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString("(" + name)
	for _, expr := range exprs {
		buf.WriteString(" ")
		s, _ := expr.Accept(ap)
		buf.WriteString(s.(string))
	}
	buf.WriteString(")")

	return buf.String(), nil
}
