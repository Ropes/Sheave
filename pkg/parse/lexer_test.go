package parse

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestLexTokens(t *testing.T) {
	var tests = []struct {
		s   string
		tok Token
		lit string
	}{
		{s: ``, tok: EOF, lit: ""},
		{s: `   `, tok: WS, lit: "   "},

		{s: `apply`, tok: CMD, lit: "apply"},
		{s: `!fish`, tok: PROMPT, lit: "fish"},
		{s: `fish`, tok: IDENT, lit: "fish"},
	}

	for i, test := range tests {
		n := fmt.Sprintf("LexTokenTest-%d-%s", i, test.s)
		t.Run(n, func(t *testing.T) {
			l, _ := NewLex(n, *bufio.NewReader(strings.NewReader(test.s)))
			tok, lit := l.Lex()

			if tok != test.tok {
				t.Errorf(" '%s' token mismatch: expect=%q got=%q %s", test.s, test.tok, tok, lit)
			}
			if lit != test.lit {
				t.Errorf("'%s' literal mismatch: expect='%s' got='%s' %v", test.s, test.lit, lit, tok)
			}
		})
	}
}
