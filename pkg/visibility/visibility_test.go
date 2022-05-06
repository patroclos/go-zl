package visibility

import (
	"fmt"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
	"jensch.works/zl/pkg/zettel/graph"
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

	g, err := graph.Make(st)
	if err != nil {
		t.Fatal(err)
	}

	visited := make([]zettel.Z, 0)
	inner := func(n crawl.Node) crawl.RecurseMask {
		visited = append(visited, n.N.Z)
		t.Log("visit", n.N.Z)
		return crawl.All
	}
	taintView := TaintView(inner, []string{"sensitive"})

	crawl.New(g, taintView).Crawl(a)

	if len(visited) != 2 {
		t.Errorf("expected 2 hits, got %d", len(visited))
	}

	visited = make([]zettel.Z, 0)
	inner = func(n crawl.Node) crawl.RecurseMask {
		visited = append(visited, n.N.Z)
		return crawl.All
	}
	taintView = TaintView(inner, nil)

	crawl.New(g, taintView).Crawl(a)

	if len(visited) != 1 {
		t.Errorf("expected 1 result, got %d", len(visited))
	}
}
