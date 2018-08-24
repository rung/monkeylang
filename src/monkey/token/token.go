package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT" // add, foobar, x, y, ...
	INT = "INT" //1,2,3,4

	ASSIGN = "="
	PLUS = "+"

	COMMA = ","
	SEMICORON = ","

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keyword
	FUNCTION = "FUNCTION"
	LET = "LET"
)



