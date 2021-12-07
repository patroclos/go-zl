package zettel

import "io"

type Id string

const InvalidId = Id("")

type Zettel interface {
	Id() Id
	Title() string
	Text() (string, error)
	SetText(t string)
	Metadata() (*MetaInfo, error)
	io.Reader
	io.Writer
	io.Seeker
}
