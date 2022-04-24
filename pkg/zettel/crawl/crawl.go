package crawl

import (
	"log"
	"strings"
	"sync"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

type RecurseMask int

const (
	None RecurseMask = 0
	Link             = LinkA | LinkB | LinkContext
	All              = Inbound | Outbound | Link
)
const (
	Inbound RecurseMask = 1 << iota
	Outbound
	LinkA
	LinkB
	LinkContext
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
	case LinkA:
		return "Link.A"
	case LinkB:
		return "Link.B"
	case LinkContext:
		return "Link.Ctx"
	}

	var parts []string
	for k := Inbound; k < LinkContext; k <<= 1 {
		if m&k != 0 {
			parts = append(parts, k.String())
		}
	}

	return strings.Join(parts, " | ")
}

type Node struct {
	Z      zettel.Z
	Path   []*Node
	Reason RecurseMask
}

type CrawlFn func(Node) RecurseMask

type Crawler interface {
	Crawl(...zettel.Z)
}

type crawlData struct {
	st zettel.ZettelerIter
	f  CrawlFn
}

func New(st zettel.ZettelerIter, f CrawlFn) Crawler {
	return crawlData{st: st, f: f}
}

func (b crawlData) Crawl(zets ...zettel.Z) {
	cr := &crawl{
		store:   b.st,
		m:       make(map[string]struct{}),
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
	store   zettel.ZettelerIter
	m       map[string]struct{}
	rw      *sync.RWMutex
	wg      *sync.WaitGroup
	errs    []error
}

func (c *crawl) Run() {
	c.errs = make([]error, 0, 8)
	c.wg.Add(len(c.root))
	for i := range c.root {
		go c.do(Node{Z: c.root[i]})
	}
	c.wg.Wait()
}

func (c *crawl) doId(id string, from Node, reason RecurseMask) {
	c.rw.Lock()
	if _, ok := c.m[id]; ok {
		c.rw.Unlock()
		return
	}
	c.rw.Unlock()
	zet, err := c.store.Zettel(id)
	if err != nil {
		c.errs = append(c.errs, err)
		log.Println(err)
		return
	}

	pth := make([]*Node, len(from.Path)+1)
	pth[len(from.Path)] = &from
	for i := range from.Path {
		pth[i] = from.Path[i]
	}
	c.wg.Add(1)
	go c.do(Node{Z: zet, Path: pth, Reason: reason})
}

func (c *crawl) do(cra Node) {
	defer c.wg.Done()
	c.rw.Lock()
	if _, ok := c.m[cra.Z.Id()]; ok {
		c.rw.Unlock()
		return
	}
	c.m[cra.Z.Id()] = struct{}{}
	c.rw.Unlock()

	mask := c.crawler(cra)

	if mask.Has(Inbound) {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			for iter := c.store.Iter(); iter.Next(); {
				boxes := scan.All(iter.Zet().Readme().Text)
				for _, box := range boxes {
					for _, ref := range box.Refs {
						zet, err := c.store.Zettel(ref[:11])
						if err != nil {
							continue
						}
						if zet.Id() == cra.Z.Id() {
							pth := make([]*Node, len(cra.Path)+1)
							pth[len(cra.Path)] = &cra
							for i := range cra.Path {
								pth[i] = cra.Path[i]
							}
							c.wg.Add(1)
							go c.do(Node{Z: iter.Zet(), Path: pth, Reason: Inbound})
						}
					}
				}
			}
		}()
	}

	if mask.Has(Outbound) {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			boxes := scan.All(cra.Z.Readme().Text)
			for _, box := range boxes {
				for _, ref := range box.Refs {
					zet, err := c.store.Zettel(ref[:11])
					if err != nil {
						continue
					}
					pth := make([]*Node, len(cra.Path)+1)
					pth[len(cra.Path)] = &cra
					for i := range cra.Path {
						pth[i] = cra.Path[i]
					}
					c.wg.Add(1)
					go c.do(Node{Z: zet, Path: pth, Reason: Outbound})
				}
			}
		}()
	}

	if lnk := cra.Z.Metadata().Link; mask&Link != None && lnk != nil {
		if mask.Has(LinkA) {
			c.doId(lnk.A, cra, LinkA)
		}
		if mask.Has(LinkB) {
			c.doId(lnk.B, cra, LinkB)
		}
		if mask.Has(LinkContext) {
			for i := range lnk.Ctx {
				c.doId(lnk.Ctx[i], cra, LinkContext)
			}
		}
	}
}
