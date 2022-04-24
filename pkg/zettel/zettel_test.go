package zettel

import (
	"strings"
	"testing"
)

func TestZetType(t *testing.T) {
	zl := &zet{}
	var _ Z = zl
}

func nullBuilder(b Builder) error { return nil }
func fooBuilder(b Builder) error {
	b.Title("foo")

	b.Text("foo\n\n* foo\n* foo\n    * foo")
	return nil
}

func TestParseReadme(t *testing.T) {
	rm := `# my cool title

my cool content`

	x, err := ParseReadme(strings.NewReader(rm))
	if err != nil {
		t.Fatal(err)
	}

	if x.Title != "my cool title" {
		t.Errorf("expected title 'my cool title', got '%s'", x.Title)
	}
	if x.Text != "my cool content" {
		t.Errorf("expected text 'my cool text', got '%s'", x.Text)
	}
}

func TestBuildErrEmptyTitle(t *testing.T) {
	zl, err := Build(nullBuilder)

	if err == nil {
		t.Error("expected error")
	}

	if zl != nil {
		t.Error("invalid Zettel returned")
	}
}

func TestBuildValid(t *testing.T) {
	got, err := Build(fooBuilder)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatalf("Zettel is nil")
	}

	if got.Id() == "" {
		t.Error("expecting non-empty result from Id()")
	}
}

func TestRebuild(t *testing.T) {
	zl, _ := Build(fooBuilder)
	id := zl.Id()
	zl2, err := zl.Rebuild(func(b Builder) error {
		b.Title("bar")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if zl2.Id() != id {
		t.Error("ids mismatch")
	}
}
