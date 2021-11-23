package graph

import (
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

func Plot(gv *graphviz.Graphviz, st storage.Storer) (*cgraph.Graph, error) {
	gr, err := gv.Graph(graphviz.Name("knowledge graph"))
	if err != nil {
		return nil, err
	}

	nodes := make(map[zettel.Id]*cgraph.Node)
	edges := make([]*cgraph.Edge, 0, 128)
	for _, z := range storage.All(st) {
		id := z.Id()
		node, err := gr.CreateNode(string(id))
		if err != nil {
			return nil, err
		}
		node.SetLabel(z.Title())
		nodes[id] = node
	}
	for _, z := range storage.All(st) {
		txt, err := z.Text()
		if err != nil {
			log.Println(err)
			continue
		}
		
		refs := zettel.Refs(txt)
		for _,ref := range refs {
			zlr,err := st.Zettel(ref)
			if err != nil {
				continue
			}
			ed,err := gr.CreateEdge(fmt.Sprintf("%s - %s", z.Id(), zlr.Id()), nodes[z.Id()], nodes[zlr.Id()])
			if err != nil {
				log.Println(err)
				return nil, err
			}
			edges = append(edges, ed)
		}
	}

	return gr, nil
}
