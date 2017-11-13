package parse

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type stateFn func(*lexer) stateFn

type item struct {
	tok Token
	val string
}

// lexer reads a stream of bytes and emits tokens
// to readers as they are discovered encoded with a
// type.
type lexer struct {
	name  string
	input bufio.Reader

	state stateFn
	items chan item
}

// NewLex returns an initialized lexer.
func NewLex(name string, input bufio.Reader) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,

		items: make(chan item),
	}
	return l, l.items
}

func (l *lexer) Lex() (tok Token, lit string) {
	ch := l.read()

	if isWhitespace(ch) {
		l.unread()
		return l.consumeWhitespace()
	} else if isBang(ch) {
		// no need to unread '!'
		return l.lexPrompt()
	} else if isAlphanumeric(ch) {
		l.unread()
		return l.lexWord()
	}

	switch ch {
	case eof:
		return EOF, ""
	}

	return ILLEGAL, string(ch)
}

func (l *lexer) emit(t Token, lit string) {
	l.items <- item{tok: t, val: lit}
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{tok: ERROR, fmt.Sprintf(format, args...)}
}

/*
func (l *lexer) run() {
	for l.state :=
}
*/

func (l *lexer) lexPrompt() (tok Token, lit string) {
	var p bytes.Buffer

	for {
		if ch := l.read(); ch == eof {
			break
		} else if !isAlphanumeric(ch) {
			l.unread()
			break
		} else {
			_, _ = p.WriteRune(ch)
		}
	}
	return PROMPT, p.String()
}

func (l *lexer) lexWord() (tok Token, lit string) {
	var w bytes.Buffer

	for {
		if ch := l.read(); ch == eof {
			break
		} else if !isAlphanumeric(ch) {
			l.unread()
			break
		} else {
			_, _ = w.WriteRune(ch)
		}
	}

	s := strings.ToLower(w.String())
	switch s {
	case "apply", "update", "list":
		return CMD, s
	}

	return IDENT, s
}

func (l *lexer) consumeWhitespace() (tok Token, lit string) {
	var ws bytes.Buffer
	ws.WriteRune(l.read()) // write the blank rune which was used to scan

	for {
		if ch := l.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			l.unread()
			break
		} else {
			ws.WriteRune(ch) // more whitespace
		}
	}
	return WS, ws.String()
}

func (l *lexer) read() rune {
	ch, _, err := l.input.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (l *lexer) unread() { _ = l.input.UnreadRune() }

// State Functions

// lexText lexically analizes a string of commands.
func lexText(l *lexer) stateFn {
	tok, lit := l.Lex()
	l.emit(tok, lit)
	switch tok {
	case PROMPT:
		return lexPrompt
	case CMD:
		return lexCmd
	default:
		return l.errorf("lexText error selecting state for literal %s", lit)
	}
}

func lexPrompt(l *lexer) stateFn {}

func lexCmd(l *lexer) stateFn {
	tok, lit := l.Lex()
	l.emit(tok, lit)

	switch tok {
	case IDENT:
		return lexIdent
	default:
		return l.errorf("lexCmd expected IDENT but got %q", tok)
	}
}

func lexIdent(l *lexer) stateFn {
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
// source: https://golang.org/src/text/template/parse/lex.go
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isBang represents a command init token
func isBang(ch rune) bool { return ch == '!' }

// From Ben Johnson's Scanner/Parser example: https://github.com/benbjohnson/sql-parser/blob/master/scanner.go
// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return unicode.IsLetter(ch) }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return unicode.IsDigit(ch) }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
