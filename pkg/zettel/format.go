package zettel

import (
	"bytes"
	"text/template"
	"time"
)

const (
	DefaultWideFormat      = `{{.Id}} {{range $k,$v := .Labels}}{{if eq $k "zl/inbox"}}📥 {{end}}{{end}} {{.Title}} {{ .Labels }}`
	ListingFormat          = `{{.Id}}  {{.Title}}`
	ListFormat             = `* {{.Id}}  {{.Title}}`
	ListStatusFormat       = `* {{.Id}}  {{.Status}} {{.Title}}`
	ListPrettyStatusFormat = `* {{.Id}}  {{.PrettyStatus}} {{.Title}}`
)

type formatData struct {
	Id           string
	Title        string
	CreateTime   time.Time
	Text         string
	Labels       map[string]string
	Inbox        *inboxData
	Link         *LinkInfo
	Status       string
	PrettyStatus string
}

type inboxData struct {
	Box string
	Due time.Time
}

func MustFmt(zl Z, format string) string {
	txt, err := Fmt(zl, format)
	if err != nil {
		panic(err)
	}

	return txt
}

func Fmt(in Z, format string) (string, error) {
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

func toFormatData(zl Z) *formatData {
	data := &formatData{
		Id:    string(zl.Id()),
		Title: zl.Readme().Title,
	}
	meta := zl.Metadata()
	data.Labels = meta.Labels
	data.Link = meta.Link
	data.CreateTime = meta.CreateTime

	if _, ok := meta.Labels["zl/inbox"]; ok {
		data.Status = "I"
		data.PrettyStatus = "📥"
	}
	return data
}
