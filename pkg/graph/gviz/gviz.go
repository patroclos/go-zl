package gviz

import (
	"fmt"
	"log"
	"sync"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

func Plot(gv *graphviz.Graphviz, st storage.Storer) (*cgraph.Graph, error) {
	gr, err := gv.Graph(graphviz.Name("knowledge graph"))
	if err != nil {
		return nil, err
	}

	taints := make(map[string]color)
	pal, ipal := []string{"red", "green", "blue", "cyan", "brown", "lime"}, 0
	for zl := range storage.AllChan(st) {
		meta, err := zl.Metadata()
		if err != nil {
			return nil, err
		}

		if val, ok := meta.Labels["zl/taint"]; ok {
			if _, ok := taints[val]; !ok {
				i := ipal
				ipal++
				col := pal[i%len(pal)]
				taints[val] = color(col)
			}
		}
	}

	nodes := make(map[zettel.Id]*cgraph.Node)
	for _, z := range storage.All(st) {
		id := z.Id()
		node, err := gr.CreateNode(string(id))
		if err != nil {
			return nil, err
		}
		node.SetLabel(z.Title())
		node.SetShape(cgraph.Box3DShape)
		nodes[id] = node

		if meta, err := z.Metadata(); err == nil {
			if taint, ok := meta.Labels["zl/taint"]; ok {
				if col, ok := taints[taint]; ok {
					node.SetFontColor(string(col))
				}
			}
		}
	}

	edges := make(chan Edge)
	wg := new(sync.WaitGroup)
	for _, z := range storage.All(st) {
		zl := z
		wg.Add(1)
		go func() {
			defer wg.Done()
			for edge := range (&EdgeScanner{
				Gr:    gr,
				St:    st,
				Nodes: nodes,
			}).Scan(zl) {
				edges <- edge
			}
		}()
	}

	go func() {
		wg.Wait()
		close(edges)
	}()

	for edge := range edges {
		_, err := gr.CreateEdge(edge.name, edge.start, edge.end)
		if err != nil {
			return nil, err
		}
	}

	wg.Wait()

	return gr, nil
}

type color string
type Edge struct {
	name  string
	start *cgraph.Node
	end   *cgraph.Node
}

func (e Edge) Edge() (name string, start *cgraph.Node, end *cgraph.Node) {
	return e.name, e.start, e.end
}

type EdgeScanner struct {
	Gr    *cgraph.Graph
	St    storage.Storer
	Nodes map[zettel.Id]*cgraph.Node
}

func (s *EdgeScanner) Scan(zl zettel.Zettel) <-chan Edge {
	ch := make(chan Edge)

	txt, err := zl.Text()
	if err != nil {
		log.Println(err)
		close(ch)
		return ch
	}

	wg := new(sync.WaitGroup)

	if meta, err := zl.Metadata(); err == nil {
		if meta.Link != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ch <- Edge{
					start: s.Nodes[zl.Id()],
					end:   s.Nodes[zettel.Id(meta.Link.A)],
					name:  string(zl.Id()),
				}
				ch <- Edge{
					start: s.Nodes[zl.Id()],
					end:   s.Nodes[zettel.Id(meta.Link.B)],
					name:  string(zl.Id()),
				}
			}()
		}
	}

	refs := scan.Refs(txt)

	wg.Add(len(refs))

	for _, ref := range refs {
		r := zettel.Id(ref)
		go func() {
			defer wg.Done()
			zlr, err := s.St.Zettel(r)
			if err != nil {
				return
			}

			name := fmt.Sprintf("%s - %s", zl.Id(), zlr.Id())
			start, end := s.Nodes[zl.Id()], s.Nodes[zlr.Id()]
			ch <- Edge{
				start: start,
				end:   end,
				name:  name,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
