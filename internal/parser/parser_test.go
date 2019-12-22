package parser

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/ozankasikci/ist/lexer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	sourceFile, err := lexer.NewSourceFile("../cmd/test.ist")
	assert.NoError(t, err)

	tokens := lexer.Lex(sourceFile)
	sourceFile.Tokens = tokens

	parseTree := Parse(sourceFile)
	assert.NotNil(t, parseTree)
	spew.Dump(parseTree)
}
