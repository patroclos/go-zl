package memory

import (
	"io"
	"time"

	"jensch.works/zl/pkg/zettel"
)

type Zettel struct {
	id      zettel.Id
	created time.Time
	title   string
	text    string
	pos     int
}

func CreateZettel(id zettel.Id, title string, text string, created time.Time) Zettel {
	return Zettel{
		id: id,
		title: title,
		text: text,
		created: created,
	}
}

func NewZettel(title string) Zettel {
	return Zettel{
		id: generateId(),
		title: title,
		text: "",
		created: time.Now(),
	}
}

func (z *Zettel) Id() zettel.Id {
	return z.id
}

func (z *Zettel) Title() string {
	return z.title
}

func (z *Zettel) CreateTime() time.Time {
	return z.created
}

func (z *Zettel) Read(p []byte) (int, error) {
	if z.pos >= len(z.text) {
		return 0, io.EOF
	}
	n := copy(p, []byte(z.text)[z.pos:])
	z.pos += n

	return n, nil
}

func (z *Zettel) Text() (string, error) {
	return z.text, nil
}

func (z *Zettel) SetText(t string) {
	z.text = t
}
