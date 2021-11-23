package zettel

import (
	"bytes"
	"text/template"
	"time"
)

const (
	DefaultWideFormat = `{{.Id}} {{range $k,$v := .Labels}}{{if eq $k "zl/inbox"}}📥 {{end}}{{end}} {{.Title}} {{ .Labels }}`
)

type formatData struct {
	Id     string
	Title  string
	CreateTime     time.Time
	Text   string
	Labels map[string]string
	Inbox  *inboxData
	Link    *LinkInfo
}

type inboxData struct {
	Box string
	Due time.Time
}

func FormatZettel(in Zettel, format string) (string, error) {
	zl := toFormatData(in)
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

func toFormatData(zl Zettel) *formatData {
	data := &formatData{
		Id: string(zl.Id()),
		Title: zl.Title(),
	}
	if meta, err := zl.Metadata(); err == nil {
		data.Labels = meta.Labels
		data.Link = meta.Link
		
		if meta.Link != nil {
		}
	}
	return data
}
