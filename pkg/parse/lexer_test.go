package parse

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
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
			c, can := context.WithCancel(context.Background())
			defer can()
			l, _ := NewLex(c, n, *bufio.NewReader(strings.NewReader(test.s)))
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

func TestLextString(t *testing.T) {
	var tests = []struct {
		s    string
		toks []Token
		lits []string
	}{
		{
			//s:    "update",
			s:    "update sheave please",
			toks: []Token{CMD, IDENT, IDENT},
			lits: []string{"update", "sheave", "please"},
		},
	}
	for i, test := range tests {
		tstr := fmt.Sprintf("lex-string-%d", i)
		t.Run(tstr, func(t *testing.T) {
			c, can := context.WithCancel(context.Background())
			defer can()
			l, _ := NewLex(c, tstr, *bufio.NewReader(strings.NewReader(test.s)))

			go l.run()

			go func() {
				i := 0
				for item := range l.items {
					t.Logf("item: %#v", item)
					if item.tok != test.toks[i] {
						t.Errorf("tokens don't match, expected %v got %v", test.toks[i], item.tok)
					}
					if item.val != test.lits[i] {
						t.Errorf("literals don't match, expected %q got %q", test.lits[i], item.val)
					}
					i++
				}
			}()
			time.Sleep(100 * time.Millisecond)
		})
	}
}
