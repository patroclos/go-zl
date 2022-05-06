package crawl

import (
	"fmt"
	"strings"
	"sync"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/elemz"
	"jensch.works/zl/pkg/zettel/graph"
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
	}

	var parts []string
	for k := Inbound; k < Outbound; k <<= 1 {
		if m&k != 0 {
			parts = append(parts, k.String())
		}
	}

	return strings.Join(parts, " | ")
}

type Node struct {
	N      *graph.Node
	Path   []*Node
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
	g graph.ZGraph
	f CrawlFn
}

func New(g graph.ZGraph, f CrawlFn) Crawler {
	return crawler{g, f}
}

func (b crawler) Crawl(zets ...zettel.Z) {
	cr := &crawl{
		g:       b.g,
		m:       make(map[int64]struct{}),
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
	g       graph.ZGraph
	m       map[int64]struct{}
	rw      *sync.RWMutex
	wg      *sync.WaitGroup
	errs    []error
}

func (c *crawl) Run() {
	c.errs = make([]error, 0, 8)
	c.wg.Add(len(c.root))
	for i := range c.root {
		go c.do(Node{
			N:      &graph.Node{Z: c.root[i]},
			Path:   []*Node{},
			Reason: Reason{},
		})
	}
	c.wg.Wait()
}

func (c *crawl) doId(id int64, from Node, reason RecurseMask) {
	c.rw.Lock()
	if _, ok := c.m[id]; ok {
		c.rw.Unlock()
		return
	}
	no := c.g.Node(id)
	if no == nil {
		c.errs = append(c.errs, fmt.Errorf("node doesnt exist %q", id))
		return
	}

	pth := make([]*Node, len(from.Path)+1)
	pth[len(from.Path)] = &from
	for i := range from.Path {
		pth[i] = from.Path[i]
	}
	c.wg.Add(1)
	go c.do(Node{N: no, Path: pth, Reason: Reason{reason, c.g.G().EdgeRefbox(from.N.ID(), id)}})
}

func (c *crawl) do(cra Node) {
	defer c.wg.Done()
	c.rw.Lock()
	id := cra.N.ID()
	if _, ok := c.m[id]; ok {
		c.rw.Unlock()
		return
	}
	c.m[id] = struct{}{}
	c.rw.Unlock()

	mask := c.crawler(cra)

	if mask.Has(Inbound) {
		inbound := c.g.G().To(id)
		for inbound.Next() {
			pth := make([]*Node, len(cra.Path)+1)
			pth[len(cra.Path)] = &cra
			for i := range cra.Path {
				pth[i] = cra.Path[i]
			}
			c.wg.Add(1)
			n := inbound.Node().(graph.Node)
			go c.do(Node{N: &n, Path: pth, Reason: Reason{Inbound, c.g.G().EdgeRefbox(inbound.Node().ID(), id)}})
		}
	}

	if mask.Has(Outbound) {
		outbound := c.g.G().From(id)
		for outbound.Next() {
			c.wg.Add(1)
			pth := make([]*Node, len(cra.Path)+1)
			pth[len(cra.Path)] = &cra
			for i := range cra.Path {
				pth[i] = cra.Path[i]
			}
			go c.do(Node{N: c.g.Node(id), Path: pth, Reason: Reason{Outbound, c.g.G().EdgeRefbox(id, outbound.Node().ID())}})
		}
	}
}
