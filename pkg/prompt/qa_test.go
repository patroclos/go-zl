package prompt

import (
	"testing"

	"git.jensch.dev/joshua/go-zl/pkg/zettel/elemz"
)

func TestParseQA(t *testing.T) {
	elemz.DefaultParser = elemz.NewCompoundParser(&parseQA{}, elemz.DefaultParser)
	elems, err := elemz.Read("Q. bla?\nA. Blub")
	if err != nil {
		t.Fatal(err)
	}
	if len(elems) != 1 {
		t.Fatalf("expected 1 *QAPrompt, got %v", elems)
	}
	switch te := elems[0].(type) {
	case *QAPrompt:
	default:
		t.Errorf("expected *QAPrompt, got %T", te)
	}

	el := elems[0]
	if el.Span().Start != 0 || el.Span().End != 15 {
		t.Errorf("expected QAPrompt.Span of <0,15>, got %v", el.Span())
	}
}
