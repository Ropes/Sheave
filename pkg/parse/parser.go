package parse

import (
	"bufio"
	"io"
)

type Command struct {
	Prompt    string
	Command   string
	Target    string
	Delimiter string
}

type Parser struct {
	l   *lexer
	buf struct {
		tok Token
		lit string
		n   int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		l: NewLexer(bufio.NewReader(r)),
	}
}
