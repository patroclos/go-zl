package zettel

import "io"

type Zettel interface {
	Id() string
	Title() string
	Metadata() *MetaInfo
	Reader() io.Reader
	Rebuild(fn func(Builder) error) (Zettel, error)
}

type zet struct {
	id    string
	title string
	meta  MetaInfo
	read  func() io.Reader
}

func (z *zet) Id() string          { return z.id }
func (z *zet) Title() string       { return z.title }
func (z *zet) Metadata() *MetaInfo { return &z.meta }
func (z *zet) Reader() io.Reader   { return z.read() }
func (z *zet) Rebuild(fn func(Builder) error) (Zettel, error) {
	b := z.toBuilder()
	if err := fn(b); err != nil {
		return nil, err
	}

	return &b.inner, nil
}

func (z *zet) toBuilder() *zettelBuilder {
	b := newBuilder()
	b.inner.id = z.id
	b.Title(z.title)
	b.inner.read = z.read
	b.inner.meta.copy(z.meta)

	return b
}
