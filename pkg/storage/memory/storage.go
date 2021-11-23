package memory

import (
	"fmt"
	"math/rand"
	"time"

	"jensch.works/zl/pkg/storage"
	z "jensch.works/zl/pkg/zettel"
)

type Storage struct {
	ZettelStorage
}

type ZettelStorage struct {
	data map[z.Id]Zettel
}

func NewStorage() storage.Storer {
	return &ZettelStorage{
		data: make(map[z.Id]Zettel),
	}
}

func (zs *ZettelStorage) HasZettel(id z.Id) bool {
	_, ok := zs.data[id]
	return ok
}

func (zs *ZettelStorage) NewZettel(title string) z.Zettel {
	zl := NewZettel(title)
	zs.data[zl.id] = zl
	return &zl
}

func (zs *ZettelStorage) SetZettel(z z.Zettel) error {
	zl, ok := z.(*Zettel)
	if !ok {
		return fmt.Errorf("Invalid zettel type. tbh this should accept anything")
	}

	zs.data[zl.id] = *zl

	return nil
}

func (zs *ZettelStorage) Zettel(id z.Id) (z.Zettel, error) {
	if val, ok := zs.data[id]; ok {
		return &val, nil
	}
	return nil, storage.ErrZettelNotFound
}

func (zs *ZettelStorage) ForEach(cb func(z.Zettel) error) error {
	for _, v := range zs.data {
		if err := cb(&v); err != nil {
			return err
		}
	}
	return nil
}

const idCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateId() z.Id {
	t := time.Now()
	rng := rand.New(rand.NewSource(t.UnixNano()))
	rbuf := [12]byte{}
	for i := 0; i < len(rbuf); i++ {
		rbuf[i] = idCharset[rng.Intn(len(idCharset))]
	}
	return z.Id(rbuf[:])
}
