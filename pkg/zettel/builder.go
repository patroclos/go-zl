package zettel

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Builder interface {
	Title(string)
	Text(string)
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
	b.Text("")

	return b
}

func (zb *zettelBuilder) Title(t string)      { zb.inner.title = t }
func (zb *zettelBuilder) Metadata() *MetaInfo { return &zb.inner.meta }

func (zb *zettelBuilder) Text(t string) {
	zb.inner.read = func() io.Reader {
		return strings.NewReader(t)
	}
}

func (zb *zettelBuilder) Validate() error {
	if len(zb.inner.id) == 0 {
		return fmt.Errorf("no id")
	}

	if len(zb.inner.title) == 0 {
		return fmt.Errorf("title empty")
	}
	if zb.inner.meta.CreateTime.IsZero() {
		return fmt.Errorf("meta.CreateTime is zero")
	}
	return nil
}
