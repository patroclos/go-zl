package zettel

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/go-git/go-billy/v5"
)

type Zettel interface {
	Id() string
	Title() string
	Metadata() *MetaInfo
	Reader() io.Reader
	Rebuild(fn func(Builder) error) (Zettel, error)
}

type Readme struct {
	Title string
	Text  string
}

func Read(id string, zd billy.Filesystem) (Zettel, error) {
	f, err := zd.Open("README.md")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	readme, err := ParseReadme(f)
	if err != nil {
		return nil, err
	}

	f, err = zd.Open("meta.yaml")
	if err != nil {
		return nil, err
	}

	meta, err := ParseMeta(f)
	if err != nil {
		return nil, err
	}

	z := &zet{
		id:    id,
		title: readme.Title,
		meta:  *meta,
		read: func() io.Reader {
			return strings.NewReader(readme.Text)
		},
	}
	return z, nil
}

func ParseReadme(r io.ReadSeeker) (*Readme, error) {
	title, err := scanTitle(r)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	txt := strings.TrimLeft(string(buf), "\n")

	return &Readme{Text: txt, Title: title}, nil
}

func scanTitle(r io.ReadSeeker) (string, error) {
	scn := bufio.NewScanner(r)
	if !scn.Scan() {
		if err := scn.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("no title")
	}

	txt := scn.Text()
	r.Seek(int64(len(txt)), io.SeekStart)
	title := strings.TrimPrefix(txt, "# ")

	return title, nil
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
