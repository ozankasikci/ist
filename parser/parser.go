package parser

import (
	"fmt"
	"github.com/ozankasikci/ist/lexer"
	"log"
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
	for p.look(0) != nil {
		if n := p.parseDecl(true); n != nil {
			p.tree.AddNode(n)	
		} else {
			log.Panicf("Unexpectec token %v")
		}
	}
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

func (p *parser) parseDecl(isTopLevel bool) ParseNode {
	var res ParseNode

    if typeDecl := p.parseTypeDecl(isTopLevel); typeDecl != nil {
    	res = typeDecl
	} else {
		return nil
	}
}

func (p *parser) parseTypeDecl(b bool) *TypeDeclNode {
    res = &TypeDeclNode{

	}
}
