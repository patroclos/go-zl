package prompt_test

import (
	"reflect"
	"testing"

	"jensch.works/zl/pkg/prompt"
)

func TestExtractQA(t *testing.T) {
	q, a := "What is the ISO 639-3 language-code for Klingon?", "tlh"
	txt := `
Q. What is the ISO 639-3 language-code for Klingon?
A. tlh
	`

	all := prompt.ExtractAll(txt)

	if len(all) != 1 {
		t.Fatalf("expected to extract 1 prompt, not %d from %s", len(all), txt)
	}

	expect := prompt.QAPrompt{
		Q: q,
		A: a,
	}

	got,ok := all[0].(prompt.QAPrompt)
	if !ok {
		t.Fatalf("Expected QAPrompt, got %#v", got)
	}

	if got.A != expect.A || got.Q != expect.Q {
		t.Errorf("expected Q and A to match: %#v, %#v", got, expect)
	}
}

// Test 
func TestExtractOmits(t *testing.T) {
	sen := "Test, {{omit}}!"
	proms := prompt.ExtractAll(sen)
	oms := make([]prompt.OmitPrompt, 0, 8)
	for _,p := range proms {
		if omit,ok := p.(prompt.OmitPrompt); ok {
			oms = append(oms, omit)
		}
	}

	l := len(oms)
	if l == 0 {
		t.Fatal("Omit not extracted")
	}

	if l > 1 {
		t.Errorf("More than one omit found for '%s': %v", sen, oms)
	}

	om := oms[0]

	if s := om.String(); s != sen {
		t.Errorf("Expected omit String() to return '%s' got '%s'", sen, s)
	}

	if l := len(om.Fragments); l != 2 {
		t.Fatalf("Expected exactly 2 fragments, got %d", l)
	}
	if l := len(om.Holes); l != 1 {
		t.Fatalf("Expected exactly 1 hole, got %d", l)
	}

	if !reflect.DeepEqual(om.Fragments, []string{"Test, ", "!"}) {
		t.Errorf("Expected Fragments to be []string{'Text, ', '!'}, found %#v", om.Fragments)
	}

	if om.Holes[0].Text != "omit" {
		t.Errorf("Expected hole to contain text 'omit', got %s", om.Holes[0].Text)
	}
}
