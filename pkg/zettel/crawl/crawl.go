package crawl

import (
	"fmt"
	"sync"

	"git.jensch.dev/joshua/zl/pkg/zettel"
	"git.jensch.dev/joshua/zl/pkg/zettel/elemz"
	"git.jensch.dev/joshua/zl/pkg/zettel/graph"
)

type RecurseMask int

const (
	None RecurseMask = 0
	All              = Inbound | Outbound
)
const (
	Inbound RecurseMask = 1 << iota
	Outbound
)

func (m RecurseMask) Has(mask RecurseMask) bool {
	return m&mask == mask
}

func (m RecurseMask) String() string {
	switch m {
	case None:
		return "None"
	case Inbound:
		return "In"
	case Outbound:
		return "Out"
	case All:
		return "All"
	default:
		return fmt.Sprintf("RecurseMask(%d)", m)
	}
}

type Node struct {
	N      graph.Node
	Path   []Node
	Reason Reason
}

type Reason struct {
	Mask   RecurseMask
	Refbox *elemz.Refbox
}

type CrawlFn func(Node) RecurseMask

type Crawler interface {
	Crawl(...zettel.Z)
}

type crawler struct {
	g *graph.Graph
	f CrawlFn
}

func New(g *graph.Graph, f CrawlFn) Crawler {
	return crawler{g, f}
}

func (b crawler) Crawl(zets ...zettel.Z) {
	cr := &crawl{
		g:       *b.g,
		seen:    make(map[int64]struct{}),
		rw:      new(sync.RWMutex),
		wg:      new(sync.WaitGroup),
		root:    zets,
		crawler: b.f,
	}
	cr.Run()
}

type crawl struct {
	root    []zettel.Z
	crawler CrawlFn
	g       graph.Graph
	seen    map[int64]struct{}
	rw      *sync.RWMutex
	wg      *sync.WaitGroup
	errs    []error
}

func (c *crawl) Run() {
	c.errs = make([]error, 0, 8)
	c.wg.Add(len(c.root))
	for i := range c.root {
		go c.do(Node{
			N:      c.g.Verts[c.g.Id(c.root[i])],
			Path:   []Node{},
			Reason: Reason{},
		})
	}
	c.wg.Wait()
}

func (c *crawl) do(cra Node) {
	defer c.wg.Done()
	c.rw.Lock()
	id := cra.N.ID()
	if _, ok := c.seen[id]; ok {
		c.rw.Unlock()
		return
	}
	c.seen[id] = struct{}{}
	c.rw.Unlock()

	mask := c.crawler(cra)

	if mask.Has(Inbound) {
		inbound := c.g.To(id)
		for inbound.Next() {
			pth := make([]Node, len(cra.Path)+1)
			pth[len(cra.Path)] = cra
			for i := range cra.Path {
				pth[i] = cra.Path[i]
			}
			n := inbound.Node().(graph.Node)
			c.wg.Add(1)
			go c.do(Node{N: n, Path: pth, Reason: Reason{Inbound, c.g.EdgeRefbox(inbound.Node().ID(), id)}})
		}
	}

	if mask.Has(Outbound) {
		outbound := c.g.From(id)
		for outbound.Next() {
			pth := make([]Node, len(cra.Path)+1)
			pth[len(cra.Path)] = cra
			for i := range cra.Path {
				pth[i] = cra.Path[i]
			}
			n := outbound.Node().(graph.Node)
			c.wg.Add(1)
			go c.do(Node{N: n, Path: pth, Reason: Reason{Outbound, c.g.EdgeRefbox(id, outbound.Node().ID())}})
		}
	}
}
