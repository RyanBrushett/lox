package glox

import (
	"errors"
	"reflect"
)

type interpreter struct{}

func NewInterpreter() *interpreter {
	return &interpreter{}
}

func (i *interpreter) visitLiteralExpr(literal *Literal) (interface{}, error) {
	return literal.Value, nil
}

func (i *interpreter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return i.evaluate(expr)
}

func (i *interpreter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	right, _ := i.evaluate(expr.Right)

	switch expr.Operator.tokenType {
	case MINUS:
		value, _ := right.(float64)
		return -value, nil
	case BANG:
		return !i.isTruthy(right), nil
	}

	return nil, nil
}

func (i *interpreter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	left, _ := i.evaluate(expr.Left)
	right, _ := i.evaluate(expr.Right)

	switch expr.Operator.tokenType {
	case MINUS:
		return left.(float64) - right.(float64), nil
	case PLUS:
		if i.isString(left) && i.isString(right) {
			return left.(string) + right.(string), nil
		} else if i.isNumeric(left) && i.isNumeric(right) {
			return left.(float64) + right.(float64), nil
		} else {
			return nil, RuntimeError(expr.Operator.line, errors.New("something's screwed up here"))
		}
	case SLASH:
		return left.(float64) / right.(float64), nil
	case STAR:
		return left.(float64) * right.(float64), nil
	case GREATER:
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case LESS:
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		return left.(float64) <= right.(float64), nil
	case EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	case BANG_EQUAL:
		return !i.isEqual(left, right), nil
	}

	return nil, nil
}

func (i *interpreter) visitTernaryExpr(expr *Ternary) (interface{}, error) {
	left, _ := i.evaluate(expr.Left)
	if i.isTruthy(left) {
		return i.evaluate(expr.Middle)
	} else {
		return i.evaluate(expr.Right)
	}
}

func (i *interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *interpreter) isTruthy(x interface{}) bool {
	if x == nil {
		return false
	}

	if reflect.TypeOf(false) == reflect.TypeOf(x) || reflect.TypeOf(true) == reflect.TypeOf(x) {
		return x.(bool)
	}

	if reflect.ValueOf(x).Kind() == reflect.Bool {
		return x.(bool)
	}

	return true
}

func (i *interpreter) isString(x interface{}) bool {
	if x == nil {
		return false
	}

	if reflect.TypeOf("") == reflect.TypeOf(x) {
		return true
	}

	if reflect.ValueOf(x).Kind() == reflect.String {
		return true
	}

	return false
}

func (i *interpreter) isNumeric(x interface{}) bool {
	if x == nil {
		return false
	}

	if reflect.TypeOf(1.0) == reflect.TypeOf(x) {
		return true
	}

	if reflect.ValueOf(x).Kind() == reflect.Float64 {
		return true
	}

	return false
}

func (i *interpreter) isEqual(x interface{}, y interface{}) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil {
		return false
	}
	if x == y {
		return true
	}
	return reflect.DeepEqual(x, y)
}
