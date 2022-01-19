package scanner

import (
	"strings"
)

type scanner struct {
	source string
}

type token struct {
	content string
}

func New(source string) *scanner {
	return &scanner{
		source: source,
	}
}

func (s *scanner) ScanTokens() []*token {
	var tokens []*token
	splitStrings := strings.Split(s.source, "\n")

	for _, s := range splitStrings {
		token := &token{content: s}
		tokens = append(tokens, token)
	}

	return tokens
}

func (t *token) String() string {
	return t.content
}
