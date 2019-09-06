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
	Where Span
}

type Position struct {
	Filename string

	Line, Char int
}

type Span struct {
	Filename string

	StartLine, StartChar int
	EndLine, EndChar int
}

func NewSpan(start, end Position) Span {
	return Span{
		Filename: start.Filename,
		StartLine: start.Line,
		StartChar: start.Char,
		EndLine: end.Line,
		EndChar: end.Char,
	}
}
