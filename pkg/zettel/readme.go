package zettel

import (
	"io"
	"strings"
)

type Readme struct {
	Title string
	Text  string
}

func ParseReadme(r io.ReadSeeker) (*Readme, error) {
	title, err := scanTitle(r)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	txt := strings.TrimLeft(string(buf), "\n")

	return &Readme{Text: txt, Title: title}, nil
}
