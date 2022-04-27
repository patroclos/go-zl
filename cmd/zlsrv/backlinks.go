package main

import (
	"strings"
	"sync"
	"time"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
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

	iter := b.s.Iter()
	for iter.Next() {
		boxes := scan.All(iter.Zet().Readme().Text)
		for _, box := range boxes {
			for _, ref := range box.Refs {
				z, err := b.s.Zettel(strings.Fields(ref)[0])
				if err != nil {
					continue
				}

				b.store(iter.Zet().Id(), z.Id())
			}
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
