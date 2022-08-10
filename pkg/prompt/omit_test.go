package prompt

import (
	"testing"

	"git.jensch.dev/zl/pkg/zettel/elemz"
)

func TestOmit(t *testing.T) {
	txt := `Something something {{hole}} something {{ another hole }}`
	ctx := &elemz.ParseCtx{
		Buf: []byte(txt),
		Pos: 0,
	}
	p := &parseOmit{}
	el, err := p.Parse(ctx)

	if err != nil {
		t.Fatal(err)
	}

	elem, ok := el.(*OmitPrompt)
	if !ok {
		t.Fatalf("expected *OmitPrompt, got %T", el)
	}

	if elem.Fragments[0] != "Something something " {
		t.Error(elem.Fragments)
	}
}
