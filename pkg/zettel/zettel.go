package zettel

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"text/template"
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

func FormatZettel(zl ZettelTemplate, format string) (string, error) {
	tmpl, err := template.New("fmt").Parse(format)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	err = tmpl.Execute(buf, zl)
	txt := string(buf.Bytes())
	if err != nil {
		return txt, err
	}
	return txt, nil
}

type ZettelTemplate struct {
	Id     string
	Title  string
	CreateTime     time.Time
	Text   string
	Labels map[string]string
	Inbox  *inboxData
	Lnk    *linkData
}

type inboxData struct {
	box string
	due time.Time
}

type linkData struct {
	a   Id   // typically the "from" end of the relationship
	b   Id   // typically the "to" end
	ctx []Id // context qualifying the relationship
}
