package zettel_test

import (
	"fmt"
	"testing"
	"time"

	"jensch.works/zl/pkg/storage/memory"
	"jensch.works/zl/pkg/zettel"
)

func makeZettel(id string, title string, labels *map[string]string) zettel.Zettel {
	zl := memory.CreateZettel(zettel.Id(id), title, "", time.Now())
	if meta, err := zl.Metadata(); err == nil && labels != nil {
		for k,v := range *labels {
			meta.Labels[k] = v
		}
	}
	return &zl
}

var (
	data1 = zettel.FormatData{
		Id: "z1",
		Title: "Zettel One",
		Labels: map[string]string{
			"zl/inbox": "default",
			"zl/taint": "work",
		},
	}
)

func ExampleFormatZettel() {
	tmpl := zettel.FormatData{
		Id:    "id1",
		Title: "Zettel 1",
		Labels: map[string]string{
			"zl/inbox": "default",
		},
	}
	msg, err := zettel.FormatZettel(tmpl, "{{.Id}} - {{.Title}}")
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)

	msg, err = zettel.FormatZettel(tmpl, `{{range $k,$v := .Labels}}{{if eq $k "zl/inbox"}}📥 {{end}}{{end}}{{.Id}} - {{.Title}}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)
	// Output:
	// id1 - Zettel 1
	// 📥 id1 - Zettel 1
}

func ExampleFormatZettelNoLabels() {
	tmpl := zettel.FormatData{
		Id:    "id1",
		Title: "Zettel 1",
	}

	msg, err := zettel.FormatZettel(tmpl, `{{range $k,$v := .Labels}}{{if eq $k "zl/inbox"}}📥 {{end}}{{end}}{{.Id}} - {{.Title}}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)
	// Output:
	// id1 - Zettel 1
}

func TestFormatWide(t *testing.T) {
	fmt := zettel.DefaultWideFormat
	expect := "z1 📥  Zettel One map[zl/inbox:default zl/taint:work]"

	txt, err := zettel.FormatZettel(data1, fmt)
	if err != nil {
		t.Fatal(err)
	}

	if txt != expect {
		t.Errorf("expected '%s', got '%s'", expect, txt)
	}
}

