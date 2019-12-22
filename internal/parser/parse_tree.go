package parser

import (
	"github.com/ozankasikci/ist/lexer"
	"math/big"
)

type ParseNode interface {
	Where() lexer.Span
	SetWhere(lexer.Span)
}

type ParseTree struct {
	Source *lexer.SourceFile
	Nodes  []ParseNode
}

func (pt ParseTree) AddNode(node ParseNode) {
	pt.Nodes = append(pt.Nodes, node)

}

type baseNode struct {
	where lexer.Span
}

func (v *baseNode) Where() lexer.Span { return v.where }
func (v *baseNode) SetWhere(w lexer.Span) { v.where = w }

type LocatedString struct {
	Where lexer.Span
	Value string
}

func NewLocatedString(token *lexer.Token) LocatedString {
	return LocatedString{Where: token.Where, Value: token.Contents}
}

type DeclNode interface {
	ParseNode
}

type baseDecl struct {
	baseNode
}

type TypeDeclNode struct {
	baseDecl
	Name LocatedString
	Type ParseNode
}

type VarDeclNode struct {
	baseDecl
	Name LocatedString
	Type *TypeReferenceNode
	Value ParseNode
}

type NameNode struct {
	baseNode
	Name    LocatedString
}

type NamedTypeNode struct {
	baseNode
	Name *NameNode
}

type TypeReferenceNode struct {
	baseNode
	Type ParseNode
}

type NumberLitNode struct {
	baseNode
	IsFloat    bool
	IntValue   *big.Int
	FloatValue float64
	FloatSize  rune
}