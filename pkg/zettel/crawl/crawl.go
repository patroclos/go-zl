package crawl

import (
	"strings"
	"sync"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

type RecurseMask int

const (
	MaskNone RecurseMask = 0
	MaskLink             = MaskLinkA | MaskLinkB | MaskLinkCtx
	MaskAll              = MaskIn | MaskOut | MaskLink
)
const (
	MaskIn RecurseMask = 1 << iota
	MaskOut
	MaskLinkA
	MaskLinkB
	MaskLinkCtx
)

func (m RecurseMask) Has(mask RecurseMask) bool {
	return m&mask == mask
}

func (m RecurseMask) String() string {
	switch m {
	case MaskNone:
		return "None"
	case MaskIn:
		return "In"
	case MaskOut:
		return "Out"
	case MaskLinkA:
		return "Link.A"
	case MaskLinkB:
		return "Link.B"
	case MaskLinkCtx:
		return "Link.Ctx"
	}

	var parts []string
	for k := MaskIn; k < MaskLinkCtx; k <<= 1 {
		if m&k != 0 {
			parts = append(parts, k.String())
		}
	}

	return strings.Join(parts, " | ")
}

type Node struct {
	Z      zettel.Zettel
	Path   []*Node
	Reason RecurseMask
}

type CrawlFn func(Node) RecurseMask

type Crawl struct {
	st zettel.ZettelerIter
	f  CrawlFn
}

func New(st zettel.ZettelerIter, f CrawlFn) Crawl {
	return Crawl{st: st, f: f}
}

func (b Crawl) Crawl(zets ...zettel.Zettel) {
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
	root    []zettel.Zettel
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

	if mask.Has(MaskIn) {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			scn := scan.ListScanner(c.store)
			for iter := c.store.Iter(); iter.Next(); {
				// TODO: cache this
				refs := scn.Scan(strings.NewReader(iter.Zet().Readme().Text))
				for ref := range refs {
					if ref.Id() == cra.Z.Id() {
						pth := make([]*Node, len(cra.Path)+1)
						pth[len(cra.Path)] = &cra
						for i := range cra.Path {
							pth[i] = cra.Path[i]
						}
						c.wg.Add(1)
						go c.do(Node{Z: iter.Zet(), Path: pth, Reason: MaskOut})
					}
				}
			}
		}()
	}

	if mask.Has(MaskOut) {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			scn := scan.ListScanner(c.store)
			for ref := range scn.Scan(strings.NewReader(cra.Z.Readme().Text)) {
				pth := make([]*Node, len(cra.Path)+1)
				pth[len(cra.Path)] = &cra
				for i := range cra.Path {
					pth[i] = cra.Path[i]
				}
				c.wg.Add(1)
				go c.do(Node{Z: ref, Path: pth, Reason: MaskOut})
			}
		}()
	}

	if lnk := cra.Z.Metadata().Link; mask&MaskLink != MaskNone && lnk != nil {
		if mask.Has(MaskLinkA) {
			zet, err := c.store.Zettel(lnk.A)
			if err != nil {
				c.errs = append(c.errs, err)
			}
			pth := make([]*Node, len(cra.Path)+1)
			pth[len(cra.Path)] = &cra
			for i := range cra.Path {
				pth[i] = cra.Path[i]
			}
			c.wg.Add(1)
			go c.do(Node{Z: zet, Path: pth, Reason: MaskLinkA})
		}
		if mask.Has(MaskLinkB) {
			zet, err := c.store.Zettel(lnk.B)
			if err != nil {
				c.errs = append(c.errs, err)
			}
			pth := make([]*Node, len(cra.Path)+1)
			pth[len(cra.Path)] = &cra
			for i := range cra.Path {
				pth[i] = cra.Path[i]
			}
			c.wg.Add(1)
			go c.do(Node{Z: zet, Path: pth, Reason: MaskLinkB})
		}
		if mask.Has(MaskLinkCtx) {
			for i := range lnk.Ctx {
				zet, err := c.store.Zettel(lnk.Ctx[i])
				if err != nil {
					c.errs = append(c.errs, err)
				}
				pth := make([]*Node, len(cra.Path)+1)
				pth[len(cra.Path)] = &cra
				for i := range cra.Path {
					pth[i] = cra.Path[i]
				}
				c.wg.Add(1)
				go c.do(Node{Z: zet, Path: pth, Reason: MaskLinkCtx})
			}
		}
	}
}
