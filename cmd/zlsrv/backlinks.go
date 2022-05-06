package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/graph"
)

type Backlinks struct {
	links map[string]map[string]struct{}
	t     time.Time
	s     zettel.ZettelerIter
	mu    sync.RWMutex
}

func (b *Backlinks) To(zet string) []string {
	blinks, ok := b.links[zet]
	if !ok {
		return nil
	}

	ids := make([]string, 0, len(blinks))
	for id := range blinks {
		ids = append(ids, id)
	}

	return ids
}

func (b *Backlinks) refresh() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.links = make(map[string]map[string]struct{})
	b.t = time.Now()

	zg, err := graph.Make(b.s)
	if err != nil {
		log.Println(fmt.Errorf("error refreshing graph: %w", err))
		return
	}

	for id := range zg.Nodes() {
		from := zg.G().From(id)
		for from.Next() {
			b.store(zg.Node(id).Z.Id(), zg.Node(from.Node().ID()).Z.Id())
		}
	}
}

func (b *Backlinks) store(from, to string) {
	if in, ok := b.links[to]; ok {
		in[from] = struct{}{}
		return
	}

	b.links[to] = map[string]struct{}{from: {}}
}
