package visibility

import (
	"fmt"
	"sync"
	"testing"

	"git.jensch.dev/joshua/zl/pkg/storage"
	"git.jensch.dev/joshua/zl/pkg/zettel"
	"git.jensch.dev/joshua/zl/pkg/zettel/crawl"
	"git.jensch.dev/joshua/zl/pkg/zettel/graph"
	"github.com/go-git/go-billy/v5/memfs"
)

// This test fails ever since moving from Nanosecond Node-IDs to UnixMilli, because
// both created Zs usually are created in the same millisecond.  The cheap fix would be to use UnixMicro or UnixNano, but
// the saner approach that I want to take would be to assign IDs to the Zs sequentially when creating the graph and storing
// those as a map
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

	if g.Edges().Len() != 2 {
		t.Fatalf("wrong edge count. expected 2, got %d", g.Edges().Len())
	}

	nodes := g.Nodes()
	if nodes.Len() != 2 {
		t.Fatal("wrong node count")
	}

	var visited []zettel.Z = nil
	var mu sync.Mutex
	inner := func(n crawl.Node) crawl.RecurseMask {
		mu.Lock()
		defer mu.Unlock()
		visited = append(visited, n.N.Z)
		return crawl.All
	}
	taintView := TaintView(inner, []string{"sensitive"})

	crawl.New(g, taintView).Crawl(a)

	if len(visited) != 2 {
		t.Errorf("expected 2 hits, got %v", visited)
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
