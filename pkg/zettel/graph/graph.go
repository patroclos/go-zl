package graph

import (
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

// https://pkg.go.dev/gonum.org/v1/gonum/graph?utm_source=godoc#Graph

type Node struct {
	Z zettel.Z
}

func (n Node) ID() int64 {
	return int64(n.Z.Metadata().CreateTime.Nanosecond())
}

func Make(store zettel.ZettelerIter) (graph.Directed, map[int64]zettel.Z, []error) {
	var errs []error

	sg := simple.NewDirectedGraph()
	idmap := make(map[int64]zettel.Z)

	iter := store.Iter()
	for iter.Next() {
		n, isNew := sg.NodeWithID(Node{iter.Zet()}.ID())
		if isNew {
			sg.AddNode(n)
			idmap[n.ID()] = iter.Zet()
		}

		boxes := scan.All(iter.Zet().Readme().Text)
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

				n2, isNew := sg.NodeWithID(Node{zet}.ID())
				if n == n2 {
					continue
				}
				if isNew {
					sg.AddNode(n2)
					idmap[n2.ID()] = zet
				}

				sg.SetEdge(sg.NewEdge(n, n2))
			}
		}
	}

	return sg, idmap, errs
}
