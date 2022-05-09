package graph

import (
	"fmt"

	"jensch.works/zl/pkg/zettel"
)

type Graph struct {
	Nodes map[int64]Node
	Zets  map[string]zettel.Z
	G     *G
}

func Id(z zettel.Z) int64 {
	return Node{z}.ID()
}

func (x Graph) Node(id int64) *Node {
	no, ok := x.Nodes[id]
	if !ok {
		return nil
	}
	return &no
}

func Make(st zettel.ZettelerIter) (*Graph, error) {
	g, idmap, err := MakeG(st)
	if len(err) > 0 {
		return nil, fmt.Errorf("%d errors: %v", len(err), err)
	}
	x := &Graph{
		Nodes: make(map[int64]Node),
		Zets:  make(map[string]zettel.Z),
		G:     g,
	}
	for id, z := range idmap {
		x.Nodes[id] = g.Node(id).(Node)
		x.Zets[z.Id()] = z
	}
	return x, nil
}
