package zettel

import (
	"strings"
	"testing"
)

func TestParseMetaLinkString(t *testing.T) {
	meta := `
link:
  from: AAA
  to: BBB`

	parsed, err := ParseMeta(strings.NewReader(meta))
	if err != nil {
		t.Fatalf("Error parsing meta: %v", err)
		return
	}

	if parsed.Link == nil {
		t.Fatalf("parsed but no link")
		return
	}

	if parsed.Link.A != "AAA" {
		t.Errorf("Expected something different than %#v from '%s'", parsed, meta)
	}
}

func TestParseMetaLinkExt(t *testing.T) {
	meta := `
link:
  from:
    zet: AAA
    aspect: something something
  to: BBB`

	parsed, err := ParseMeta(strings.NewReader(meta))
	if err != nil {
		t.Fatalf("Error parsing meta: %v", err)
		return
	}

	if parsed.Link == nil {
		t.Fatalf("parsed but no link")
		return
	}

	if err := validateLink(parsed.Link); err != nil {
		t.Errorf("(failed validation: %v)", err)
	}

	if parsed.Link.A != "AAA" {
		t.Errorf("expected link.from to be AAA but got %s", parsed.Link.A)
	}
	// FIXME the aspect gets thrown away here
}
