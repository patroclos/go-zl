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
	io.Reader
}

func Refs(zettel Zettel) []Id {
	results := []Id{}
	buf := new(strings.Builder)
	_, err := io.Copy(buf, zettel)
	if err != nil {
		return make([]Id, 0)
	}
	reg := regexp.MustCompile(`\[.*?\]\((.*?)\)`)
	matches := reg.FindAllStringSubmatch(buf.String(), 0)
	for _, m := range matches {
		results = append(results, Id(m[1]))
	}

	return results
}
