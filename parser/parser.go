package parser

import (
	"github.com/ozankasikci/ist/lexer"
)

type parser struct {
	source       *lexer.SourceFile
	currentToken int
	tree         *ParseTree
}

func Parse(input *lexer.SourceFile) (*ParseTree) {
	p := &parser{}

    p.parse()

	return p.tree
}

func (p *parser) parse()  {

}

func (p *parser) look(ahead int) *lexer.Token {
	if ahead < 0 {
		panic("look method can not accept a negative value!")
	}

	if p.currentToken + ahead >= len(p.source.Tokens) {
		return nil
	}

	return p.source.Tokens[p.currentToken + ahead]
}
