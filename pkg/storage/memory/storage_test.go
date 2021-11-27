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

func TestZettelIO(t *testing.T) {
	store := memory.NewStorage()
	zettel := store.NewZettel("blub")

	text := "Hello, Zettel!"
	zettel.SetText(text)

	buf := make([]byte, 64)
	read, err := zettel.Read(buf)
	if err != nil {
		t.Fatalf("Error reading back: %v", err)
	}
	if read != 14 {
		t.Fatalf("Expected to read back 14 bytes. got %d", read)
	}
}

func BenchmarkStorageInsert(b *testing.B) {
	st := memory.NewStorage()
	for i := 0; i < b.N; i++ {
		st.NewZettel("test")
	}
}
