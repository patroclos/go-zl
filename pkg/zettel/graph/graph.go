package graph

import (
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/elemz"
)

// G is a gonum-powered Z-graph providing access to parsed refbox.
// This struct wraps a simple graph.Directed implementation, bolting on metadata
// about boxes, relations, etc.
type Graph struct {
	*simple.DirectedGraph
	Verts map[int64]Node
	Zets  map[string]zettel.Z
	boxes map[int64]map[int64]elemz.Refbox
}

func (g *Graph) Refbox(from, to zettel.Z) (elemz.Refbox, bool) {
	x, ok := g.boxes[Node{from}.ID()][Node{to}.ID()]
	return x, ok
}

func Id(z zettel.Z) int64 {
	return Node{z}.ID()
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
	}

	iter := st.Iter()
	for iter.Next() {
		n := Node{iter.Zet()}
		nId := n.ID()

		if _, ok := g.Verts[nId]; !ok {
			g.AddNode(n)
			g.Verts[nId] = n
			g.Zets[n.Z.Id()] = n.Z
		}

		bm := g.boxes[nId]
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
				id := strings.Fields(ref)[0]

				zet, err := st.Zettel(id)
				if err != nil {
					continue
				}

				n2 := Node{Z: zet}
				if nId == n2.ID() {
					continue
				}

				if _, ok := g.Verts[n2.ID()]; !ok {
					g.AddNode(n2)
					g.Verts[n2.ID()] = n2
					g.Zets[n2.Z.Id()] = n2.Z
				}

				bm[n2.ID()] = box
				g.SetEdge(g.NewEdge(n, n2))
			}
		} // end range boxes

		g.boxes[nId] = bm
	} // end iter.Next
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
