package context

import (
	"jensch.works/zl/pkg/zettel"
)

type Context struct {
	Store    zettel.Storage
	Template string
	Labels   []zettel.Labelspec
}
