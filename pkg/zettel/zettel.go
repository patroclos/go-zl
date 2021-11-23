package zettel

import (
	"io"
	"regexp"
	"strings"
	"time"
)

type Id string

type Zettel interface {
	Id() Id
	Title() string
	CreateTime() time.Time
	Text() (string, error)
	SetText(t string)
	Metadata() (*MetaInfo, error)
	io.Reader
	io.Writer
	io.Seeker
}

func Refs(text string) []Id {
	reg := regexp.MustCompile(`\[.+\]\((.+)\)`)
	matches := reg.FindAllStringSubmatch(text, -1)
	results := make([]Id, 0, 8)
	for _, m := range matches {
		id := strings.Trim(m[1], " /")
		results = append(results, Id(id))
	}
	reg = regexp.MustCompile(`\* ([a-zA-Z0-9-]+)  .*`)
	matches = reg.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		id := m[1]
		results = append(results, Id(id))
	}

	return results
}

