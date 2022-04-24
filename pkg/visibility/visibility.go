package visibility

import (
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
)

func TaintView(inner crawl.CrawlFn, tolerate []string) crawl.CrawlFn {
	return func(n crawl.Node) crawl.RecurseMask {
		if !Visible(n.Z, tolerate) {
			return crawl.None
		}
		return inner(n)
	}
}

func Visible(z zettel.Z, tolerate []string) bool {
	if len(tolerate) == 1 && tolerate[0] == "*" {
		return true
	}
	taint, ok := z.Metadata().Labels["zl/taint"]
	if !ok {
		return true
	}

	for i := range tolerate {
		if taint == tolerate[i] {
			return true
		}
	}
	return false
}
