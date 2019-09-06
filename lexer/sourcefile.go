package lexer

import (
	"io/ioutil"
	"path/filepath"
)

type SourceFile struct {
	Path, Name string
	Contents []rune
	NewLines []int
	Tokens []*Token
}

func NewSourceFile(path string) (*SourceFile, error)  {
	name := filepath.Base(path)

	source := &SourceFile{Name: name, Path: path}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	source.Contents = []rune(string(contents))
	return source, nil
}
