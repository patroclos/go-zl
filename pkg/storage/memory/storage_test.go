package memory_test

import (
	"testing"

	"jensch.works/zl/pkg/storage/memory"
)

func TestZettelCreation(t *testing.T) {
	store := memory.NewStorage()
	title := "TestZettelCreation"
	zettel := store.NewZettel(title)

	if got := zettel.Title(); got != title {
		t.Fatalf("expected zettel.Title() to return '%s', but got '%s'", title, got)
	}
}
