package main

import (
	"github.com/ozankasikci/ist/lexer"
	"github.com/y0ssar1an/q"
	"log"
)

func main() {
	ifile, err := lexer.NewInputFile("/Users/ozankasikci/Documents/projects/ist/cmd/test.ist")
	if err != nil {
		log.Fatal(err)
	}

    tokens := lexer.Lex(ifile)
    q.Q(tokens)
}
