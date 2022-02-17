package visibility

import "jensch.works/zl/pkg/zettel/crawl"

func TaintView(inner crawl.CrawlFn, tolerate []string) crawl.CrawlFn {
	return func(n crawl.Node) crawl.RecurseMask {
		taint, ok := n.Z.Metadata().Labels["zl/taint"]
		if !ok {
			return inner(n)
		}

		for i := range tolerate {
			if tolerate[i] == taint {
				return inner(n)
			}
		}

		return crawl.None
	}
}
