package lexer

import (
	"fmt"
	"github.com/y0ssar1an/q"
	"log"
	"unicode"
)

type lexer struct {
	input            *inputFile
	startPos, endPos int
	currentPos       Position
	tokenStart       Position
}

func (l *lexer) lex() {
	for {
		l.skipLayoutAndComments()

		if isEOF(l.lookAhead(0)) {
			l.input.NewLines = append(l.input.NewLines, l.endPos)
			return
		} else if isLetter(l.lookAhead(0)) || l.lookAhead(0) == '_' {
			l.recognizeIdentifierToken()
		} else {
			log.Panicf("unrecognized token")
		}
	}
}

func Lex(i *inputFile) []*Token {
	l := &lexer{
		input:      i,
		startPos:   0,
		endPos:     0,
		currentPos: Position{Filename: i.Name, Line: 1, Char: 1},
		tokenStart: Position{Filename: i.Name, Line: 1, Char: 1},
	}

	l.lex()

	return l.input.Tokens
}

func (l *lexer) lookAhead(distance int) rune {
	if distance < 0 {
		panic(fmt.Sprintf("Tried to lookAhead a negative number: %d", distance))
	}

	if l.endPos+distance >= len(l.input.Contents) {
		return 0
	}
	return l.input.Contents[l.endPos+distance]
}

func (l*lexer) recognizeIdentifierToken() {
	l.consume()

	for isLetter(l.lookAhead(0)) || isDecimalDigit(l.lookAhead(0)) || l.lookAhead(0) == '_' {
		l.consume()
	}

	l.pushToken(Identifier)
}

func (l *lexer) consume()  {
	l.currentPos.Char += 1

	if isEOL(l.lookAhead(0)) {
		l.currentPos.Char = 1
		l.currentPos.Line += 1
		l.input.NewLines = append(l.input.NewLines, l.endPos)
	}

	l.endPos += 1
}

func (l*lexer) pushToken(t TokenType) {
	tok := &Token{
		Type:     t,
		Contents: string(l.input.Contents[l.startPos:l.endPos]),
	}

	l.input.Tokens = append(l.input.Tokens, tok)
	q.Q("lexer", "[%4d:%4d:% 11s] `%s`\n", l.startPos, l.endPos, tok.Type, tok.Contents)
	l.flushBuffer()
}

func (l *lexer) flushBuffer() {
	l.startPos = l.endPos

	l.tokenStart = l.currentPos
}

func (l *lexer) skipLayoutAndComments() {
	for {
		for isLayout(l.lookAhead(0)) {
			l.consume()
		}
		l.flushBuffer()
		break
	}

	//v.printBuffer()
	l.flushBuffer()
}

func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isEOL(r rune) bool {
	return r == '\n'
}

func isEOF(r rune) bool {
	return r == 0
}

// IsSpace reports whether the rune is a space character as defined
// by Unicode's White Space property; in the Latin-1 space
// this is
//	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
func isLayout(r rune) bool {
	return (r <= ' ' || unicode.IsSpace(r)) && !isEOF(r)
}
