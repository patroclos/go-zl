package zettel

import "io"

type Zettel interface {
	Id() string
	Title() string
	Metadata() *MetaInfo
	Reader() io.Reader
}
