package zettel

import (
	"fmt"
	"io"
	"os"
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

func (rm Readme) String() string {
	return fmt.Sprintf("# %s\n\n%s", rm.Title, rm.Text)
}

func (rm Readme) NewTemp() (*os.File, error) {
	tmp, err := os.CreateTemp("", "zledit*.md")
	if err != nil {
		return nil, err
	}

	tmp.Write([]byte(rm.String()))

	return tmp, nil
}
