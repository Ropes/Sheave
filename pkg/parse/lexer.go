package parse

import (
	"bufio"
	"bytes"
	"strings"
)

type stateFn func(*lexer) stateFn

type lexer struct {
	name  string
	input bufio.Reader
}

func NewLex(name string, input bufio.Reader) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
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

// isAlphanumeric returns true if rune is in the alphabet or digit.
func isAlphanumeric(ch rune) bool {
	if isLetter(ch) || isDigit(ch) {
		return true
	}
	return false
}

func isBang(ch rune) bool { return ch == '!' }

// From Ben Johnson's Scanner/Parser example: https://github.com/benbjohnson/sql-parser/blob/master/scanner.go
// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
