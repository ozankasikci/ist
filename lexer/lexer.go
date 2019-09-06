package lexer

import (
	"fmt"
	"github.com/y0ssar1an/q"
	"log"
	"os"
	"strings"
	"unicode"
)

type lexer struct {
	source *SourceFile

	// the position in the source.Contents rune slice
	bufferStart, bufferEnd int

	//  the position in the SourceFile source
	currentPos, tokenStart Position
}

func (l *lexer) lex() {
	for {
		// skips new lines and comments
		l.skipSpaceAndComments()

		// the end of file, return
		if isEOF(l.look(0)) {
			l.source.NewLines = append(l.source.NewLines, l.bufferEnd)
			return
		}

		// the rest must be a token
		if isLetter(l.look(0)) || l.look(0) == '_' {
			l.recognizeIdentifierToken()
		} else if l.look(0) == '"' {
			l.recognizeStringToken()
		} else if isNumber(l.look(0)) {
			l.recognizeNumberToken()
		} else if isOperator(l.look(0)) {
			l.recognizeOperatorToken()
		} else {
			log.Panicf("unrecognized token")
		}
	}
}

func Lex(i *SourceFile) []*Token {
	l := &lexer{
		source:      i,
		bufferStart: 0,
		bufferEnd:   0,
		currentPos:  Position{Filename: i.Name, Line: 1, Char: 1},
		tokenStart:  Position{Filename: i.Name, Line: 1, Char: 1},
	}

	l.lex()

	return l.source.Tokens
}

func (l *lexer) look(distance int) rune {
	if distance < 0 {
		panic(fmt.Sprintf("Tried to look a negative number: %d", distance))
	}

	if l.bufferEnd+distance >= len(l.source.Contents) {
		return 0
	}
	return l.source.Contents[l.bufferEnd+distance]
}

func (l *lexer) recognizeIdentifierToken() {
	l.consume()

	for isLetter(l.look(0)) || isDecimalDigit(l.look(0)) || l.look(0) == '_' {
		l.consume()
	}

	l.pushToken(Identifier)
}

func (l *lexer) recognizeNumberToken() {
	l.consume()

	for isNumber(l.look(0)) {
		l.consume()
	}

	l.pushToken(Number)
}

func (l *lexer) recognizeStringToken() {
	pos := l.currentPos

	// prepare to read string value
	l.flushBuffer()

	for {
		if l.look(0) == '"' {
			// end of string, push token
			l.pushToken(String)
			l.consume()
			return
		} else if isEOF(l.look(0)) {
			// end of file without ending string literal, exit
			l.errPos(pos, "Unterminated string literal!")
		} else {
			l.consume()
		}
	}
}

func (l *lexer) recognizeOperatorToken() {
	pos := l.currentPos

	if strings.ContainsRune("=!<>", l.look(0)) && l.look(1) == '=' {
		l.consume()
		l.consume()
	} else {
		l.errPos(pos, "Unexpected operator!")
	}

	l.pushToken(Operator)
}

func (l *lexer) consume() {
	l.currentPos.Char += 1

	if isEOL(l.look(0)) {
		l.currentPos.Char = 1
		l.currentPos.Line += 1
		l.source.NewLines = append(l.source.NewLines, l.bufferEnd)
	}

	l.bufferEnd += 1
}

func (l *lexer) pushToken(t TokenType) {
	tok := &Token{
		Type:     t,
		Contents: string(l.source.Contents[l.bufferStart:l.bufferEnd]),
	}

	l.source.Tokens = append(l.source.Tokens, tok)
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

func (l *lexer) errPos(pos Position, err string, rest ...interface{}) {
	log.Printf("Position: %v", pos)
	println(err)

	os.Exit(1)
}

func isDecimalDigit(r rune) bool { return r >= '0' && r <= '9' }
func isLetter(r rune) bool       { return unicode.IsLetter(r) }
func isEOL(r rune) bool          { return r == '\n' }
func isEOF(r rune) bool          { return r == 0 }
func isNumber(r rune) bool       { return unicode.IsNumber(r) }
func isOperator(r rune) bool     { return strings.ContainsRune("+-*/=><!", r) }
func isSpace(r rune) bool {
	// IsSpace reports whether the rune is a space character as defined
	// by Unicode's White Space property; in the Latin-1 space
	// this is
	//	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
	return (r <= ' ' || unicode.IsSpace(r)) && !isEOF(r)
}
