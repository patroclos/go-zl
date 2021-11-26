package storage

import (
	"errors"
	"io"

	z "jensch.works/zl/pkg/zettel"
)

var (
	ErrZettelNotFound = errors.New("zettel not found")
)

type Storer interface {
	ZettelStorer
	ZettelIter
}

type ZettelStorer interface {
	Zetteler
	NewZettel(title string) z.Zettel
	SetZettel(zettel z.Zettel) error
	HasZettel(id z.Id) bool
}

type Zetteler interface {
	Zettel(id z.Id) (z.Zettel, error)
}

type emptyZettler struct{}

var Empty = emptyZettler{}

func (_ emptyZettler) Zettel(id z.Id) (z.Zettel, error) {
	return nil, io.EOF
}

type ZettelIter interface {
	ForEach(func(z.Zettel) error) error
}

func All(iter ZettelIter) []z.Zettel {
	results := make([]z.Zettel, 0, 512)
	iter.ForEach(func(z z.Zettel) error {
		results = append(results, z)
		return nil
	})

	return results
}

func AllChan(iter ZettelIter) <-chan z.Zettel {
	ch := make(chan z.Zettel)
	go func() {
		defer close(ch)
		iter.ForEach(func(zl z.Zettel) error{
			ch <- zl
			return nil
		})
	}()

	return ch
}
