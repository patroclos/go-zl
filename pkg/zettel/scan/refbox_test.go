package scan

import (
	"testing"
)

const pattern = `(?m)^([a-zA-Z][^:\n]*):\n((?:\* [0-9]{6}-[a-zA-Z]{4}  .+$\n)+)(\+ .*$(?:\n  .*$)*)?`

func TestRefboxString(t *testing.T) {
	box := Refbox{
		Rel:   "Related",
		Refs:  []string{"060102-blub  Hello Zet"},
		Extra: []string{"somehow this works", "savy?"},
	}

	str := box.String()

	if str != `Related:
* 060102-blub  Hello Zet
+ somehow this works
  savy?` {
		t.Errorf("expectation failed, got %q", str)
	}
}

func TestRefboxAll(t *testing.T) {
	txt := `Heiß und Fettig:
* 220122-NiZh  First title
* 220122-x8oH  Another one
+ Freeform extra text
  Of course with multiline support 🤯`

	boxes := All(txt)

	if len(boxes) != 1 {
		t.Fatalf("expected to find 1 Refbox, got %d", len(boxes))
	}

	b := boxes[0]
	if b.Start != 0 || b.End != len(txt) {
		t.Errorf("expected start,end to be 0,%d, got %d,%d", len(txt), b.Start, b.End)
	}

	if len(b.Refs) != 2 {
		t.Fatalf("Expected 2 refs, got %d", len(b.Refs))
	}

	if len(b.Extra) != 2 {
		t.Fatalf("Expected 2 extra lines, got %d", len(b.Extra))
	}

	for i, x := range []string{"220122-NiZh  First title", "220122-x8oH  Another one"} {
		if b.Refs[i] != x {
			t.Errorf("expected Refs[%d] to be %q, got %q", i, x, b.Refs[i])
		}
	}

	for i, x := range []string{"Freeform extra text", "Of course with multiline support 🤯"} {
		if b.Extra[i] != x {
			t.Errorf("expected Refs[%d] to be %q, got %q", i, x, b.Extra[i])
		}
	}

	if str := b.String(); str != txt {
		t.Errorf("expected refbox.String to match original input, got %#v", str)
	}
}

func TestRefboxNoTrailingNewline(t *testing.T) {
	txt := "Refs:\n* 220101-blub  TITLE"
	all := All(txt)
	if l := len(all); l != 1 {
		t.Errorf("expected to parse 1 refbox from %q, got %d", txt, l)
	}
}

func TestRefboxUri(t *testing.T) {
	txt := "Refs:\n* <https://jensch.dev>"
	all := All(txt)
	if l := len(all); l != 1 {
		t.Errorf("expected %q to contain a refbox, got %d", txt, l)
	}
}
