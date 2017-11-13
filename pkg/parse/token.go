package parse

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS
	ERROR

	// Prompt to execute
	PROMPT
	// Command or action to take
	CMD

	IDENT // IDENTifying variable
)
