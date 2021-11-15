package storage

import (
	"errors"

	z "jensch.works/zl/pkg/zettel"
)

var (
	ErrZettelNotFound = errors.New("zettel not found")
)

type Storer interface {
	ZettelStorer
}

type ZettelStorer interface {
	NewZettel(title string) z.Zettel
	SetZettel(zettel z.Zettel) error
	Zettel(id z.Id) (z.Zettel, error)
	IterZettel() (ZettelIter, error)
}

type ZettelIter interface {
	ForEach(func(z.Zettel) error) error
}


