package zettel_test

import (
	"fmt"
	"testing"

	"jensch.works/zl/pkg/storage/memory"
	"jensch.works/zl/pkg/zettel"
)

func makeZettel(id string, title string, labels *map[string]string) zettel.Zettel {
	zl := memory.CreateZettel(zettel.Id(id), title, "")
	if meta, err := zl.Metadata(); err == nil && labels != nil {
		*meta = zettel.MetaInfo{
			Labels: *labels,
		}
	}
	return &zl
}

var (
	testZettel = makeZettel("z1", "Zettel One", &map[string]string{
		"zl/inbox": "default",
		"zl/taint": "work",
	})
)

func ExampleFormatZettel() {
	zl := makeZettel("id1", "Zettel 1", &map[string]string{"zl/inbox": "default"})
	msg, err := zettel.FormatZettel(zl, "{{.Id}} - {{.Title}}")
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)

	msg, err = zettel.FormatZettel(zl, zettel.DefaultWideFormat)
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)
	// Output:
	// id1 - Zettel 1
	// id1 📥  Zettel 1 map[zl/inbox:default]
}


func TestFormatWide(t *testing.T) {
	fmt := zettel.DefaultWideFormat
	expect := "z1 📥  Zettel One map[zl/inbox:default zl/taint:work]"

	txt, err := zettel.FormatZettel(testZettel, fmt)
	if err != nil {
		t.Fatal(err)
	}

	if txt != expect {
		t.Errorf("expected '%s', got '%s'", expect, txt)
	}
}
