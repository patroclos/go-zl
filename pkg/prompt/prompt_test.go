package prompt_test

import (
	"testing"

	"jensch.works/zl/pkg/prompt"
)

func TestExtractPrompts(t *testing.T) {
	q, a := "What is the ISO 639-3 language-code for Klingon?", "tlh"
	txt := `
Q. What is the ISO 639-3 language-code for Klingon?
A. tlh
	`

	got := prompt.ExtractAll(txt)

	if len(got) != 1 {
		t.Fatalf("expected to extract 1 prompt, not %d from %s", len(got), txt)
	}

	expect := prompt.EmbeddedPrompt{
		Q: q,
		A: a,
	}

	if got[0].A != expect.A || got[0].Q != expect.Q {
		t.Errorf("expected Q and A to match: %#v, %#v", got[0], expect)
	}
}
