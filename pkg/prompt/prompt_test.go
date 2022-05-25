package prompt

import (
	"testing"

	"jensch.works/zl/pkg/zettel/elemz"
)

func TestExtractQA(t *testing.T) {
	q, a := "What is the ISO 639-3 language-code for Klingon?", "tlh"
	txt := `Q. What is the ISO 639-3 language-code for Klingon?
A. tlh`

	all, err := elemz.ReadWith(txt, &elemz.OneOfParser{Parsers: []elemz.Parser{&parseQA{}}})
	if err != nil {
		t.Fatalf("failed parsing QA: %v", err)
	}

	if len(all) != 1 {
		t.Fatalf("expected to extract 1 prompt, not %d from %s", len(all), txt)
	}

	expect := QAPrompt{
		Q: q,
		A: a,
	}

	got, ok := all[0].(*QAPrompt)
	if !ok {
		t.Fatalf("Expected *QAPrompt, got %T", all[0])
	}

	if got.A != expect.A || got.Q != expect.Q {
		t.Errorf("expected Q and A to match: %#v, %#v", got, expect)
	}
}

// Test
func TestExtractOmits(t *testing.T) {
	sen := "Test, {{omit}}!"
	proms, err := elemz.ReadWith(sen, &parseOmit{})
	if err != nil {
		t.Fatal(err)
	}

	if len(proms) != 1 {
		t.Fatal("invalid result cound", len(proms))
	}
	if proms[0].ElemType() != OmitType {
		t.Fatalf("expected omit, got %v", proms[0].ElemType())
	}
}
