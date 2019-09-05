package lexer

import (
	"io/ioutil"
	"path/filepath"
)

type inputFile struct {
	Path, Name string
	Contents []rune
	NewLines []int
	Tokens []*Token
}

func NewInputFile(path string) (*inputFile, error)  {
	name := filepath.Base(path)

	iFile := &inputFile{Name: name, Path: path}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	iFile.Contents = []rune(string(contents))
	return iFile, nil
}
