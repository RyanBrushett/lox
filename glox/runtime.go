package glox

import (
	"fmt"
)

type runtime struct {
	HadError bool
}

type runtimeError struct {
	line  int
	where int
	Err   error
}

func RuntimeError(line int, err error) *runtimeError {
	return &runtimeError{line, 0, err}
}

func (e *runtimeError) Error() string {
	return fmt.Sprintf("[Line %d] Error%d: %v\n", e.line, e.where, e.Err)
}

func NewRuntime() *runtime {
	return &runtime{
		HadError: false,
	}
}

func (r *runtime) Run(source string) {
	tokens, err := NewScanner(source).ScanTokens()

	if err != nil {
		r.reportError(err)
	}

	// for now just print
	for _, t := range tokens {
		fmt.Printf("=> %s\n", t)
	}
}

func (r *runtime) reportError(e error) {
	fmt.Println(e)
	r.HadError = true
}
