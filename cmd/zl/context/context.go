package context

import (
	"jensch.works/zl/pkg/zettel"
)

type Context struct {
	Template string
	Labels   []zettel.Labelspec
}
