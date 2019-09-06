package parser

import "github.com/ozankasikci/ist/lexer"

type ParseNode interface {
	Where() lexer.Span
}

type ParseTree struct {
	Source *lexer.SourceFile
	Nodes []ParseNode
}

type baseNode struct {
	where lexer.Span
}

func (v *baseNode) Where() lexer.Span                { return v.where }

type LocatedString struct {
	Where lexer.Span
	Value string
}

func NewLocatedString(token *lexer.Token) LocatedString {
	return LocatedString{ Where: token.Where, Value: token.Contents }
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

func (pt ParseTree) AddNode(node ParseNode)  {
	pt.Nodes = append(pt.Nodes, node)

}
