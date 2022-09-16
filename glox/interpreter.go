package glox

import (
	"errors"
	"fmt"
	"reflect"
)

type interpreter struct{}

func NewInterpreter() *interpreter {
	return &interpreter{}
}

func (i *interpreter) Interpret(statements []Stmt, inREPL bool) error {
	for _, statement := range statements {
		value, err := i.execute(statement)
		if err != nil {
			return err
		}

		if inREPL && value != nil {
			fmt.Printf("%v\n", value)
		}
	}
	return nil
}

func (i *interpreter) execute(statement Stmt) (interface{}, error) {
	return statement.Accept(i)
}

func (i *interpreter) visitVarStmt(stmt *Var) (interface{}, error) {
	return i.evaluate(stmt.Initializer) // TODO
}

func (i *interpreter) visitExpressionStmt(stmt *Expression) (interface{}, error) {
	return i.evaluate(stmt.Expression)
}

func (i *interpreter) visitPrintStmt(stmt *Print) (interface{}, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", value)
	return nil, nil
}

func (i *interpreter) visitLiteralExpr(literal *Literal) (interface{}, error) {
	return literal.Value, nil
}

func (i *interpreter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return i.evaluate(expr.Expression)
}

func (i *interpreter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	right, _ := i.evaluate(expr.Right)

	switch expr.Operator.tokenType {
	case MINUS:
		err := i.checkNumericOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
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
		err := i.checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case PLUS:
		if i.isString(left) && i.isString(right) {
			return left.(string) + right.(string), nil
		} else if i.isNumeric(left) && i.isNumeric(right) {
			return left.(float64) + right.(float64), nil
		} else {
			return nil, errors.New("operands in addition must both be numeric or both be strings")
		}
	case SLASH:
		if err := i.checkNumericOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		if err := i.checkDivideByZero(right.(float64)); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case STAR:
		err := i.checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case GREATER:
		err := i.checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		err := i.checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		err := i.checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		err := i.checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
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

func (i *interpreter) visitVariableExpr(expr *Variable) (interface{}, error) {
	return nil, nil // TODO
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

func (i *interpreter) checkNumericOperand(operator *Token, x interface{}) error {
	if _, ok := x.(float64); !ok {
		return fmt.Errorf(
			"operand '%v' in '%s' operation is not a numeric value",
			x, operator.lexeme,
		)
	}
	return nil
}

func (i *interpreter) checkNumericOperands(operator *Token, left, right interface{}) error {
	for _, x := range []interface{}{left, right} {
		if err := i.checkNumericOperand(operator, x); err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) checkDivideByZero(right float64) error {
	if right == 0 {
		return fmt.Errorf("cannot divide by zero")
	}
	return nil
}
