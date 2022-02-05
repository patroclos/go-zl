package zettel

import (
	"fmt"
	"time"
)

type Builder interface {
	Id(string)
	Title(string)
	Text(string)
	Metadata() *MetaInfo
}

func Build(fn func(b Builder) error) (Zettel, error) {
	b := newBuilder()

	if err := fn(b); err != nil {
		return nil, err
	}
	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate builder: %w", err)
	}
	return &b.inner, nil
}

type zettelBuilder struct {
	inner zet
}

func newBuilder() *zettelBuilder {
	b := &zettelBuilder{}
	b.inner.meta.CreateTime = time.Now()
	b.inner.id = plainGenerateId()

	return b
}

func (zb *zettelBuilder) Id(id string) {
	zb.inner.id = id
}

func (zb *zettelBuilder) Title(t string) {
	zb.inner.readme.Title = t
}

func (zb *zettelBuilder) Text(t string) {
	zb.inner.readme.Text = t
}

func (zb *zettelBuilder) Metadata() *MetaInfo {
	return &zb.inner.meta
}

func (zb *zettelBuilder) Validate() error {
	if len(zb.inner.id) == 0 {
		return fmt.Errorf("no id")
	}

	if len(zb.inner.readme.Title) == 0 {
		return fmt.Errorf("title empty")
	}
	if zb.inner.meta.CreateTime.IsZero() {
		return fmt.Errorf("meta.CreateTime is zero")
	}
	return nil
}
