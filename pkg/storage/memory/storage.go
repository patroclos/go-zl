package memory

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"jensch.works/zl/pkg/storage"
	z "jensch.works/zl/pkg/zettel"
)

type Storage struct {
	ZettelStorage
}

type ZettelStorage struct {
	data map[z.ZId]zettel
}

func NewStorage() storage.Storer {
	return &ZettelStorage{
		data: make(map[z.ZId]zettel),
	}
}

func (zs *ZettelStorage) NewZettel(title string) z.Zettel {
	zl := zettel{
		id:      generateId(),
		title:   title,
		content: "",
		created: time.Now(),
	}
	zs.data[zl.id] = zl
	return &zl
}

func (zs *ZettelStorage) SetZettel(z z.Zettel) error {
	zl, ok := z.(*zettel)
	if !ok {
		return fmt.Errorf("Invalid zettel type")
	}

	zs.data[zl.id] = *zl

	return nil
}

func (zs *ZettelStorage) Zettel(id z.ZId) (z.Zettel, error) {
	if val, ok := zs.data[id]; ok {
		return &val, nil
	}
	return nil, storage.ErrZettelNotFound
}

func (zs *ZettelStorage) IterZettel() (storage.ZettelIter, error) {
	i := iter{
		zs: zs,
	}
	return &i, nil
}

type iter struct {
	zs *ZettelStorage
}

func (i *iter) ForEach(cb func(z.Zettel) error) error {
	for _, v := range i.zs.data {
		if err := cb(&v); err != nil {
			return err
		}
	}
	return nil
}

type zettel struct {
	id      z.ZId
	created time.Time
	title   string
	content string
}

func (z *zettel) Id() z.ZId {
	return z.id
}

func (z *zettel) Title() string {
	return z.title
}

func (z *zettel) CreateTime() time.Time {
	return z.created
}

func (z *zettel) Reader() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(z.content)), nil
}

const idCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateId() z.ZId {
	t := time.Now()
	// y, m, d := t.Date()
	// h, min, s := t.Clock()
	rng := rand.New(rand.NewSource(t.UnixNano()))
	rbuf := [12]byte{}
	for i := 0; i < len(rbuf); i++ {
		rbuf[i] = idCharset[rng.Intn(len(idCharset))]
	}
	return z.ZId(rbuf[:])
	// return z.ZId(fmt.Sprintf("%02d%02d%02d-%02d%02d%02d-%s", y, m, d, h, min, s, rbuf))
}
