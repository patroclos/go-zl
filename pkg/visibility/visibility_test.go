package visibility

import (
	"fmt"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
)

func TestTaintView(t *testing.T) {
	st, _ := storage.NewStore(memfs.New())
	id1, id2 := zettel.MakeId(), zettel.MakeId()
	a, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(id1)
		b.Title("one")
		b.Text(fmt.Sprintf("Refs:\n* %s  Two", id2))
		return nil
	})
	b, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(id2)
		b.Title("Two")
		b.Text(fmt.Sprintf("Refs:\n* %s  One", id1))
		b.Metadata().Labels["zl/taint"] = "sensitive"
		return nil
	})
	st.Put(a)
	st.Put(b)

	visited := make([]zettel.Zettel, 0)
	inner := func(n crawl.Node) crawl.RecurseMask {
		visited = append(visited, n.Z)
		t.Log("visit", n.Z)
		return crawl.All
	}
	taintView := TaintView(inner, []string{"sensitive"})

	crawl.New(st, taintView).Crawl(a)

	if len(visited) != 2 {
		t.Errorf("expected 2 hits, got %d", len(visited))
	}

	visited = make([]zettel.Zettel, 0)
	inner = func(n crawl.Node) crawl.RecurseMask {
		visited = append(visited, n.Z)
		return crawl.All
	}
	taintView = TaintView(inner, nil)

	crawl.New(st, taintView).Crawl(a)

	if len(visited) != 1 {
		t.Errorf("expected 1 result, got %d", len(visited))
	}
}
