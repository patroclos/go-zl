package zettel_test

import (
	"testing"

	"jensch.works/zl/pkg/zettel"
)

func TestTemplate(t *testing.T) {
	fmt := "{{ .Id }}  {{ .Title }}"
	dat := zettel.ZettelTemplate{
		Id:    "testid",
		Title: "Testerino",
	}

	txt, err := zettel.FormatZettel(dat, fmt)
	if err != nil {
		t.Error(err)
	}

	expected := "testid  Testerino"
	if txt != expected {
		t.Errorf("Expected '%s', got '%s'", expected, txt)
	}

}

func TestRefParsing(t *testing.T) {
	in := `
* id1  Title One
* blub single space doesnt match
Embedded [Second](id2/)
`
	refs := zettel.Refs(in)

	expectation := []zettel.Id{zettel.Id("id1"), zettel.Id("id2")}
	for _, i := range expectation {
		found := false
		for _, r := range refs {
			if string(r) == string(i) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find ref to %s in %s", i, in)
		}
	}
}
