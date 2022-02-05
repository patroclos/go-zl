package zettel

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/go-git/go-billy/v5"
	"gopkg.in/yaml.v2"
)

type Zettel interface {
	Id() string
	Readme() Readme
	Metadata() *MetaInfo
	Rebuild(fn func(Builder) error) (Zettel, error)
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
		id:     id,
		readme: *readme,
		meta:   *meta,
	}
	return z, nil
}

func Write(zet Zettel, dir billy.Filesystem) error {
	if err := writeReadme(zet, dir); err != nil {
		return err
	}
	if err := writeMeta(zet, dir); err != nil {
		return err
	}
	return nil
}

func writeReadme(zet Zettel, dir billy.Filesystem) error {
	readme := zet.Readme()
	fReadme, err := dir.Create("README.md")
	if err != nil {
		return err
	}
	defer fReadme.Close()

	_, err = fmt.Fprintf(fReadme, "%s", readme.String())

	if err != nil {
		return err
	}

	return nil
}

func writeMeta(zet Zettel, dir billy.Filesystem) error {
	meta := zet.Metadata()
	fMeta, err := dir.Create("meta.yaml")
	if err != nil {
		return err
	}
	defer fMeta.Close()

	mb, err := yaml.Marshal(meta)
	if err != nil {
		return err
	}

	_, err = io.Copy(fMeta, bytes.NewReader(mb))
	if err != nil {
		return err
	}

	return nil
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
	id     string
	readme Readme
	meta   MetaInfo
}

func (z *zet) Id() string          { return z.id }
func (z *zet) Readme() Readme      { return z.readme }
func (z *zet) Metadata() *MetaInfo { return &z.meta }
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
	b.inner.readme.Title = z.readme.Title
	b.inner.readme.Text = z.readme.Text
	b.inner.meta.copy(z.meta)

	return b
}
