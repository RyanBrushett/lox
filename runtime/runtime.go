package runtime

import (
	"fmt"
	"lox/scanner"
)

type runtime struct {
	HadError bool
}

func New() *runtime {
	return &runtime{
		HadError: false,
	}
}

func (r *runtime) Run(source string) {
	tokens := scanner.New(source).ScanTokens()

	for _, t := range tokens {
		fmt.Printf("=> %s\n", t)
	}
}

func (r *runtime) Error(line int, message string) {
	r.report(line, "", message)
}

func (r *runtime) report(line int, where string, message string) {
	fmt.Printf("[Line %d] Error%s: %s\n", line, where, message)
	r.HadError = true
}
