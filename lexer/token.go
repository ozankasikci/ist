package lexer

type TokenType int

const (
	Rune TokenType = iota
	Identifier
	Operator
	Number
	String
)

func (tt TokenType) String() string  {
	return [...]string{"rune", "identifier", "operator", "number", "string"}[tt]
}

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

func (s Span) Start() Position {
	return Position{
		Filename: s.Filename,
		Line: s.StartLine,
		Char: s.StartChar,
	}
}

func (s Span) End() Position {
	return Position{
		Filename: s.Filename,
		Line: s.EndLine,
		Char: s.EndChar,
	}
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
