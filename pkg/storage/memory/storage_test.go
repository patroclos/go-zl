package memory_test

import (
	"fmt"
	"testing"

	"jensch.works/zl/pkg/storage/memory"
	"jensch.works/zl/pkg/zettel"
)

func TestZettelCreation(t *testing.T) {
	store := memory.NewStorage()
	title := "TestZettelCreation"
	zettel := store.NewZettel(title)

	if got := zettel.Title(); got != title {
		t.Fatalf("expected zettel.Title() to return '%s', but got '%s'", title, got)
	}
}

func TestZettelIO(t *testing.T) {
	store := memory.NewStorage()
	zettel := store.NewZettel("blub")

	text := "Hello, Zettel!"
	zettel.SetText(text)

	buf := make([]byte, 64)
	read, err := zettel.Read(buf)
	if err != nil {
		t.Fatalf("Error reading back: %v", err)
	}
	if read != 14 {
		t.Fatalf("Expected to read back 14 bytes. got %d", read)
	}
}

func ExampleFormatZettel() {
	tmpl := zettel.ZettelTemplate{
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
	tmpl := zettel.ZettelTemplate{
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
