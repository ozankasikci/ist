package parser

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/ozankasikci/ist/lexer"
	"log"
	"math/big"
	"strings"
)

type parser struct {
	source       *lexer.SourceFile
	currentToken int
	tree         *ParseTree
}

func Parse(source *lexer.SourceFile) (*ParseTree) {
	p := &parser{
		source: source,
		tree: &ParseTree{Source: source},
	}

    p.parse()

	return p.tree
}

func (p *parser) parse()  {
	for p.look(0) != nil {
		if n := p.parseDecl(true); n != nil {
			println("parsed decl")
			p.tree.AddNode(n)
		} else {
			log.Panicf("Unexpected token %v", p.look(0))
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
    	println("parseTypeDecl 1")
		res = typeDecl
	} else if varDecl := p.parseVarDecl(isTopLevel); varDecl != nil {
		println("parseVarDecl 1")
		res = varDecl
	} else {
		return nil
	}

    return res
}

func (p *parser) parseTypeDecl(isTopLevel bool) *TypeDeclNode {
    return nil
}

func (p *parser) parseVarDecl(isTopLevel bool) *VarDeclNode {
	body := p.parseVarDeclBody()
	if body == nil {
		return nil
	}

	return body
}

func (p *parser) parseVarDeclBody() *VarDeclNode {

	if !p.tokensMatch(lexer.Identifier, "", lexer.Identifier, "") {
		println("tokens not match")
		return nil
	}

	varName := p.consumeToken()

	varType := p.parseTypeReference(true)
	if varType == nil && !p.tokenMatches(0, lexer.Operator, "=") {
		p.err("Expected valid type in variable declaration")
	}

	var value ParseNode
	if p.tokenMatches(0, lexer.Operator, "=") {
		p.consumeToken()

		value = p.parseExpr()

		if value == nil {
			p.err("parseVarDeclBody value is null")
		}
	}

    res := &VarDeclNode{
    	Name: NewLocatedString(varName),
    	Type: varType,
	}

    start := varName.Where.Start()

    var end lexer.Position
    if value != nil {
    	res.Value = value
    	end = value.Where().End()
	} else {
		end = varType.Where().End()
	}

    res.SetWhere(lexer.NewSpan(start, end))
	return res
}

func (p *parser) parseType(mustParse bool) ParseNode {
    var res ParseNode

	if p.nextTokenIs(lexer.Identifier) {
		res = p.parseNamedType()
	} else {
		println("ParseType else")
	}

    return res
}

func (p *parser) parseTypeReference(mustParse bool) *TypeReferenceNode {
	typ := p.parseType(mustParse)
	if typ == nil {
		return nil
	}

	res := &TypeReferenceNode{
		Type: typ,
	}

	res.SetWhere(lexer.NewSpan(typ.Where().Start(), typ.Where().End()))

	return res
}

func (p *parser) parseNamedType() *NamedTypeNode {
	name := p.parseName()
	if name == nil {
		return nil
	}

	res := &NamedTypeNode{Name: name}
	res.SetWhere(name.Where())
	return res
}

func (p *parser) parseName() *NameNode {

	if !p.nextTokenIs(lexer.Identifier)  {
		println("next is not identifier")
		return nil
	}

	name := p.consumeToken()
    res := &NameNode{ Name: NewLocatedString(name) }
    res.SetWhere(name.Where)

    return res
}

func (p *parser) parseExpr() ParseNode {
	pri := p.parsePostfixExpr()
	if pri == nil {
		return  nil
	}

	return pri
}

func (p *parser) parsePostfixExpr() ParseNode {
	expr := p.parsePrimaryExpr()
	if expr == nil {
		return nil
	}

	return expr
}

func (p *parser) parsePrimaryExpr() ParseNode {
	var res ParseNode

	if litExpr := p.parseLitExpr(); litExpr != nil {
		res = litExpr
	}

	return res
}

func (p *parser) parseLitExpr() ParseNode {
	var res ParseNode

	if numberLiteral := p.parseNumberLiteral(); numberLiteral != nil {
		res = numberLiteral
	}

	return res
}

func (p *parser) parseNumberLiteral() ParseNode {
	if !p.nextTokenIs(lexer.Number) {
		return nil
	}

	t:= p.consumeToken()
	number := t.Contents
	res := &NumberLitNode{}

	intValue, ok := parseInt(number, 10)
	if !ok {
		p.err("cant parse int")
	}

	res.IntValue = intValue

	res.SetWhere(t.Where)
	return res
}

func (p *parser) tokenMatches(ahead int, t lexer.TokenType, contents string) bool {
	token := p.look(ahead)
	return token != nil && token.Type == t && (contents == "" || token.Contents == contents)
}

func (p *parser) consumeToken() *lexer.Token {
	ret := p.look(0)
	p.currentToken += 1
	return ret
}

func (p *parser) nextTokenIs(typ lexer.TokenType) bool {
	next := p.look(0)
	if next == nil {
		p.err("")
		log.Panicf("Expected token of type %s, got EOF", typ)
	}

	return next.Type == typ
}

func (p *parser) expect(typ lexer.TokenType, val string) *lexer.Token {
	if !p.tokenMatches(0, typ, val) {
		// token doesnt match
		t := p.look(0)
		if t == nil {
			if val != "" {
				log.Panicf("Expected `%s` (%s), got EOF", val, typ)
			} else {
				log.Panicf("Expected %s, got EOF", typ)
			}
		} else {
			if val != "" {
				log.Panicf("Expected `%s` (%s), got `%s` (%s)", val, typ, t.Contents, t.Type)
			} else {
				log.Panicf("Expected %s, got %s (`%s`)", typ, t.Type, t.Contents)
			}
		}
	}

	return p.consumeToken()
}

func (p *parser) err(text string)  {
	println("Parser error: ", text)
	spew.Dump(struct{
		CurrentToken int
		CurrentTokenValue *lexer.Token
		TokenCount int
	}{
		p.currentToken,
		p.look(0),
		len(p.source.Tokens),
	})
}

func parseInt(num string, base int) (*big.Int, bool) {
	num = strings.ToLower(strings.Replace(num, "_", "", -1))

	var splitNum []string
	if base == 10 {
		splitNum = strings.Split(num, "e")
	} else {
		splitNum = []string{num}
	}

	if !(len(splitNum) == 1 || len(splitNum) == 2) {
		return nil, false
	}

	numVal := splitNum[0]

	ret := big.NewInt(0)

	_, ok := ret.SetString(numVal, base)
	if !ok {
		return nil, false
	}

	// handle standard form
	if len(splitNum) == 2 {
		expVal := splitNum[1]

		exp := big.NewInt(0)
		_, ok = exp.SetString(expVal, base)
		if !ok {
			return nil, false
		}

		if exp.BitLen() > 64 {
			panic("TODO handle this better")
		}
		expInt := exp.Int64()

		ten := big.NewInt(10)

		if expInt < 0 {
			for ; expInt < 0; expInt++ {
				ret.Div(ret, ten)
			}
		} else if expInt > 0 {
			for ; expInt > 0; expInt-- {
				ret.Mul(ret, ten)
			}
		}
	}

	return ret, true
}

func (p*parser) tokensMatch(args ...interface{}) bool {
	if len(args) % 2 != 0 {
		panic("passed uneven args to tokensMatch")
	}

	for i := 0; i < len(args) / 2; i++ {
		if !(p.tokenMatches(i, args[i * 2].(lexer.TokenType), args[i * 2 + 1].(string))) {
			return false
		}
	}
	return true
}
