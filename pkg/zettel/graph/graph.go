package graph

import (
	"strings"

	"gonum.org/v1/gonum/graph/simple"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/elemz"
)

// https://pkg.go.dev/gonum.org/v1/gonum/graph?utm_source=godoc#Graph

type Node struct {
	Z zettel.Z
}

func (n Node) ID() int64 {
	return int64(n.Z.Metadata().CreateTime.Nanosecond())
}

type G struct {
	*simple.DirectedGraph
	boxes map[int64]map[int64]elemz.Refbox
}

func (g *G) EdgeRefbox(from, to int64) *elemz.Refbox {
	if !g.HasEdgeFromTo(from, to) {
		return nil
	}
	boxes, ok := g.boxes[from]
	if !ok {
		return nil
	}
	box, ok := boxes[to]
	if !ok {
		return nil
	}
	return &box
}

func MakeG(store zettel.ZettelerIter) (*G, map[int64]zettel.Z, []error) {
	var errs []error

	sg := simple.NewDirectedGraph()
	idmap := make(map[int64]zettel.Z)

	boxmap := map[int64]map[int64]elemz.Refbox{}
	iter := store.Iter()
	for iter.Next() {
		n := Node{iter.Zet()}
		if _, ok := idmap[n.ID()]; !ok {
			sg.AddNode(n)
			idmap[n.ID()] = n.Z
		}

		bm := boxmap[n.ID()]
		if bm == nil {
			bm = make(map[int64]elemz.Refbox)
		}
		boxes := elemz.Refboxes(iter.Zet().Readme().Text)
		for _, box := range boxes {
			for _, ref := range box.Refs {
				id := strings.Fields(ref)[0]
				if strings.HasPrefix(id, "<") {
					continue
				}
				zet, err := store.Zettel(id)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				n2 := Node{zet}
				if n.ID() == n2.ID() {
					continue
				}
				if _, ok := idmap[n2.ID()]; !ok {
					sg.AddNode(n2)
					idmap[n2.ID()] = zet
				}
				bm[n2.ID()] = box

				sg.SetEdge(sg.NewEdge(n, n2))
			}
		}
		boxmap[n.ID()] = bm
	}

	return &G{sg, boxmap}, idmap, errs
}
