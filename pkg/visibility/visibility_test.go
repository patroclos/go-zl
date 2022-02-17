package visibility

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
)

func TestTaintView(t *testing.T) {
	st, _ := storage.NewStore(memfs.New())
	zl, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id("one")
		b.Title("one")
		b.Text("Refs:\n* two  Two")
		return nil
	})
	st.Put(zl)
	zl, _ = zettel.Build(func(b zettel.Builder) error {
		b.Id("two")
		b.Title("Two")
		b.Text("Refs:\n* one  One")
		b.Metadata().Labels["zl/taint"] = "sensitive"
		return nil
	})
	st.Put(zl)

	visited := make([]zettel.Zettel, 0)
	inner := func(n crawl.Node) crawl.RecurseMask {
		visited = append(visited, n.Z)
		return crawl.All
	}
	taintView := TaintView(inner, []string{"sensitive"})

	zl, _ = st.Zettel("one")
	crawl.New(st, taintView).Crawl(zl)

	if len(visited) != 2 {
		t.Errorf("expected 2 hits, got %d", len(visited))
	}

	visited = make([]zettel.Zettel, 0)
	taintView = TaintView(inner, nil)

	crawl.New(st, taintView).Crawl(zl)

	if len(visited) != 1 {
		t.Errorf("expected 1 result, got %d", len(visited))
	}
}
