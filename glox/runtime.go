package glox

import (
	"fmt"
	"os"
)

type loxRuntime struct {
	HadError bool
}

type parseError struct {
	line  int
	where string
	Err   error
}

type runtimeError struct {
	line  int
	where string
	Err   error
}

func RuntimeError(line int, err error) *runtimeError {
	return &runtimeError{line, "0", err}
}

func ParseError(token *Token, err error) *parseError {
	if token.tokenType == EOF {
		return &parseError{token.line, " at end", err}
	} else {
		at := fmt.Sprintf(" at '%s'", token.lexeme)
		return &parseError{token.line, at, err}
	}
}

func (e *parseError) Error() string {
	return fmt.Sprintf("[Line %d] Error%s: %v\n", e.line, e.where, e.Err)
}

func (e *runtimeError) Error() string {
	return fmt.Sprintf("[Line %d] Error%s: %v\n", e.line, e.where, e.Err)
}

func NewRuntime() *loxRuntime {
	return &loxRuntime{
		HadError: false,
	}
}

func (r *loxRuntime) Run(source string, line int) {
	if source == "exit" || source == "exit!" || source == "quit" {
		os.Exit(0)
	}

	tokens, err := NewScanner(source).ScanTokens()

	if err != nil {
		r.reportError(err)
		return
	}

	parser := NewParser(tokens)
	statements := parser.parse()
	interpreter := NewInterpreter()
	err = interpreter.Interpret(statements)

	if err != nil {
		r.reportError(RuntimeError(line, err))
		return
	}
}

func (r *loxRuntime) reportError(e error) {
	fmt.Println(e)
	r.HadError = true
}
