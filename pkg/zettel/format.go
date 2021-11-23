package zettel

import (
	"bytes"
	"text/template"
	"time"
)

const (
	DefaultWideFormat = `{{.Id}} {{range $k,$v := .Labels}}{{if eq $k "zl/inbox"}}📥 {{end}}{{end}} {{.Title}} {{ .Labels }}`
)


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

func toFormatData(zl Zettel) *FormatData {
	data := &FormatData{
		Id: string(zl.Id()),
		Title: zl.Title(),
	}
	if meta, err := zl.Metadata(); err == nil {
		data.Labels = meta.Labels
	}
	return data
}

type FormatData struct {
	Id     string
	Title  string
	CreateTime     time.Time
	Text   string
	Labels map[string]string
	Inbox  *InboxData
	Lnk    *LinkData
}

type InboxData struct {
	Box string
	Due time.Time
}

type LinkData struct {
	A   Id   // typically the "from" end of the relationship
	B   Id   // typically the "to" end
	Ctx []Id // context qualifying the relationship
}
