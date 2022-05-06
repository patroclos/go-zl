package elemz

import (
	"testing"
)

func TestElems(t *testing.T) {
	txt := "```mermaid\nblub\n```\n\nblarb"

	x, err := Read(txt)
	if err != nil {
		t.Fatal(err)
	}

	if len(x) != 2 {
		t.Fatalf("expected 2 elemz, got %d (%v)", len(x), x)
	}

	c, ok := x[0].(*Code)
	if !ok {
		t.Fatalf("expected *Code, got %T", x[0])
	}

	if c.BlockParam != "mermaid" {
		t.Errorf("Code-Block parameter expected mermaid, got %q", c.BlockParam)
	}

	if c.Code != "blub" {
		t.Errorf("code blub expected, got %q", c.Code)
	}

	xt, ok := x[1].(*Text)
	if !ok {
		t.Fatalf("expected *Text, got %T", x[1])
	}
	expect := "\n\nblarb"
	if xt.String() != expect {
		t.Errorf("expected text %q, got %q", expect, xt.String())
	}
}
