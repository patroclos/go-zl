package zettel

import "io"

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
