package memory

import (
	"bytes"
	"errors"
	"io"
	"time"

	"jensch.works/zl/pkg/zettel"
)

type Zettel struct {
	id      zettel.Id
	created time.Time
	title   string
	buf     *bytes.Buffer
	pos     int
	meta    zettel.MetaInfo
}

func CreateZettel(id zettel.Id, title string, text string, created time.Time) Zettel {
	return Zettel{
		id:      id,
		title:   title,
		buf:     bytes.NewBufferString(text),
		created: created,
	}
}

func NewZettel(title string) Zettel {
	return Zettel{
		id:      generateId(),
		title:   title,
		buf:     bytes.NewBuffer(nil),
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
	if z.pos == z.buf.Len() {
		return 0, io.EOF
	}
	n := copy(p, z.buf.Bytes()[z.pos:])
	z.pos += n
	return n, nil
}

func (z *Zettel) Write(p []byte) (n int, err error) {
	if z.pos < z.buf.Len() {
		n = copy(z.buf.Bytes()[z.pos:], p)
		p = p[n:]
	}

	if len(p) > 0 {
		var bn int
		bn, err = z.buf.Write(p)
		n += bn
	}

	z.pos += n
	return n, err
}

func (z *Zettel) Seek(off int64, whence int) (int64, error) {
	newPos, offs := 0, int(off)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = z.pos + offs
	case io.SeekEnd:
		newPos = z.buf.Len() + offs
	}

	if newPos < 0 {
		return 0, errors.New("seek before start")
	}
	if newPos > z.buf.Len() {
		newPos = z.buf.Len()
	}
	z.pos = newPos
	return int64(newPos), nil
}

func (z *Zettel) Reader() io.Reader {
	return bytes.NewReader(z.buf.Bytes())
}

func (z *Zettel) Text() (string, error) {
	return z.buf.String(), nil
}

func (z *Zettel) SetText(t string) {
	z.buf = bytes.NewBufferString(t)
	z.pos = 0
}

func (z *Zettel) Metadata() (*zettel.MetaInfo, error) {
	return &z.meta, nil
}
