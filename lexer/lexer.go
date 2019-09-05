package lexer

import (
	"fmt"
	"github.com/y0ssar1an/q"
	"log"
	"unicode"
)

type lexer struct {
	input            *inputFile
	// the position in the input.Contents rune slice
	bufferStart, bufferEnd int
	currentPos             Position
	tokenStart             Position
}

func (l *lexer) lex() {
	for {
		// skips new lines and comments
		l.skipSpaceAndComments()

		// the end of file, return
		if isEOF(l.look(0)) {
			l.input.NewLines = append(l.input.NewLines, l.bufferEnd)
			return
		}

		// the rest must be a token
		if isLetter(l.look(0)) || l.look(0) == '_' {
			l.recognizeIdentifierToken()
		} else if isNumber(l.look(0)) {
			l.recognizeNumberToken()
		} else {
			log.Panicf("unrecognized token")
		}
	}
}

func Lex(i *inputFile) []*Token {
	l := &lexer{
		input:       i,
		bufferStart: 0,
		bufferEnd:   0,
		currentPos:  Position{Filename: i.Name, Line: 1, Char: 1},
		tokenStart:  Position{Filename: i.Name, Line: 1, Char: 1},
	}

	l.lex()

	return l.input.Tokens
}

func (l *lexer) look(distance int) rune {
	if distance < 0 {
		panic(fmt.Sprintf("Tried to look a negative number: %d", distance))
	}

	if l.bufferEnd+distance >= len(l.input.Contents) {
		return 0
	}
	return l.input.Contents[l.bufferEnd+distance]
}

func (l*lexer) recognizeIdentifierToken() {
	l.consume()

	for isLetter(l.look(0)) || isDecimalDigit(l.look(0)) || l.look(0) == '_' {
		l.consume()
	}

	l.pushToken(Identifier)
}

func (l *lexer) recognizeNumberToken()  {
	l.consume()
	l.pushToken(Number)
}

func (l *lexer) consume()  {
	l.currentPos.Char += 1

	if isEOL(l.look(0)) {
		l.currentPos.Char = 1
		l.currentPos.Line += 1
		l.input.NewLines = append(l.input.NewLines, l.bufferEnd)
	}

	l.bufferEnd += 1
}

func (l*lexer) pushToken(t TokenType) {
	tok := &Token{
		Type:     t,
		Contents: string(l.input.Contents[l.bufferStart:l.bufferEnd]),
	}

	l.input.Tokens = append(l.input.Tokens, tok)
	q.Q("lexer", "[%4d:%4d:% 11s] `%s`\n", l.bufferStart, l.bufferEnd, tok.Type, tok.Contents)
	l.flushBuffer()
}

func (l *lexer) flushBuffer() {
	l.bufferStart = l.bufferEnd

	l.tokenStart = l.currentPos
}

func (l *lexer) skipSpaceAndComments() {
	for {
		for isSpace(l.look(0)) {
			l.consume()
		}
		l.flushBuffer()
		break
	}

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

func isNumber(r rune) bool {
	return unicode.IsNumber(r)
}

func isSpace(r rune) bool {
	// IsSpace reports whether the rune is a space character as defined
	// by Unicode's White Space property; in the Latin-1 space
	// this is
	//	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
	return (r <= ' ' || unicode.IsSpace(r)) && !isEOF(r)
}
