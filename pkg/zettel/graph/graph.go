package graph

import (
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"git.jensch.dev/joshua/go-zl/pkg/zettel"
	"git.jensch.dev/joshua/go-zl/pkg/zettel/elemz"
)

type Node struct {
	Z  zettel.Z
	id int64
}

func (n Node) ID() int64 {
	return n.id
}

// G is a gonum-powered Z-graph providing access to parsed refbox.
// This struct wraps a simple graph.Directed implementation, bolting on metadata
// about boxes, relations, etc.
type Graph struct {
	*simple.DirectedGraph
	Verts map[int64]Node
	Zets  map[string]zettel.Z
	ids   map[string]int64
	boxes map[int64]map[int64]elemz.Refbox
}

func (g *Graph) Refbox(from, to zettel.Z) (elemz.Refbox, bool) {
	x, ok := g.boxes[g.Id(from)][g.Id(to)]
	return x, ok
}

func (g *Graph) Id(z zettel.Z) int64 {
	return g.ids[z.Id()]
}

func (g *Graph) NodeZ(id int64) *Node {
	no, ok := g.Verts[id]
	if !ok {
		return nil
	}
	return &no
}

func (g *Graph) Node(id int64) graph.Node {
	return g.NodeZ(id)
}

func Make(st zettel.ZettelerIter) (*Graph, error) {
	g := &Graph{
		DirectedGraph: simple.NewDirectedGraph(),
		Verts:         map[int64]Node{},
		Zets:          map[string]zettel.Z{},
		boxes:         map[int64]map[int64]elemz.Refbox{},
		ids:           map[string]int64{},
	}

	iter := st.Iter()
	var nextId int64
	for iter.Next() {
		nextId++
		n := Node{Z: iter.Zet(), id: nextId}
		g.AddNode(n)
		g.Verts[n.ID()] = n
		g.Zets[n.Z.Id()] = n.Z
		g.ids[n.Z.Id()] = n.ID()
	}
	for _, n := range g.Verts {
		bm := g.boxes[n.ID()]
		if bm == nil {
			bm = map[int64]elemz.Refbox{}
		}

		boxes := elemz.Refboxes(n.Z.Readme().Text)
		for _, box := range boxes {
			for _, ref := range box.Refs {
				if strings.HasPrefix(ref, "<") {
					// TODO: emit an egress knode
					continue
				}
				id, ok := g.ids[strings.Fields(ref)[0]]
				if !ok || n.id == id {
					continue
				}

				n2 := g.Verts[id]
				bm[n2.ID()] = box
				g.SetEdge(g.NewEdge(n, n2))
			}
		}
		g.boxes[n.ID()] = bm
	}
	return g, nil
}

func (g *Graph) EdgeRefbox(from, to int64) *elemz.Refbox {
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

/*
func MakeG(store zettel.ZettelerIter) (*Graph, map[int64]zettel.Z, []error) {
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

	return &Graph{
		DirectedGraph: 5,
		sg: boxmap,
	}, idmap, errs
}
*/
