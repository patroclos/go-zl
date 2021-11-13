package zettel

import (
	"io"
	"regexp"
	"strings"
	"time"
)

type ZId string

type Zettel interface {
	Id() ZId
	Title() string
	CreateTime() time.Time
	Reader() (io.ReadCloser, error)
}

func Refs(zettel Zettel) []ZId {
	results := []ZId{}
	buf := new(strings.Builder)
	reader, err := zettel.Reader()
	if err != nil {
		return make([]ZId, 0)
	}
	_, err = io.Copy(buf, reader)
	if err != nil {
		return make([]ZId, 0)
	}
	reg := regexp.MustCompile(`\[.*?\]\((.*?)\)`)
	matches := reg.FindAllStringSubmatch(buf.String(), 0)
	for _, m := range matches {
		results = append(results, ZId(m[1]))
	}

	return results
}

