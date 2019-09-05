package lexer

type TokenType int

const (
	Rune TokenType = iota
	Identifier
	Number
	String
	Operator
)

type Token struct {
	Type     TokenType
	Contents string
}

type Position struct {
	Filename string

	Line, Char int
}
