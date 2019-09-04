package lexer

type lexer struct {
	input *inputFile
	startPos, endPos int
	currentPos Position
	tokenStart  Position

}

func (l *lexer) lex()  {
	
}

func Lex(i *inputFile) []*Token {
	return nil
}
