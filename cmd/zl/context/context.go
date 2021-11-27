package context

import (
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

type Context struct {
	Template string
	Store    storage.Storer
	Labels   []zettel.Labelspec
}
