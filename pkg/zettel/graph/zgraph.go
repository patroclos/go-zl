package graph

import (
	"fmt"

	"jensch.works/zl/pkg/zettel"
)

type ZGraph interface {
	Node(id int64) *Node
	Nodes() map[int64]Node
	G() *G
}

type zg struct {
	nodes map[int64]Node
	zets  map[string]zettel.Z
	g     *G
}

func (x zg) Nodes() map[int64]Node {
	return x.nodes
}

func (x zg) Node(id int64) *Node {
	no, ok := x.nodes[id]
	if !ok {
		return nil
	}
	return &no
}

func (x zg) G() *G {
	return x.g
}

func Make(st zettel.ZettelerIter) (ZGraph, error) {
	g, idmap, err := MakeG(st)
	if len(err) > 0 {
		return nil, fmt.Errorf("%d errors: %v", len(err), err)
	}
	x := zg{
		nodes: make(map[int64]Node),
		zets:  make(map[string]zettel.Z),
		g:     g,
	}
	for id, z := range idmap {
		x.nodes[id] = g.Node(id).(Node)
		x.zets[z.Id()] = z
	}
	return x, nil
}
