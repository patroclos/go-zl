package elemz

import (
	"fmt"
	"io"
)

type ElemType string

type Elem interface {
	Spaner
	fmt.Stringer
	ElemType() ElemType
}

type Spaner interface {
	Span() Span
}

type Span struct {
	Start int
	End   int
}

func Read(txt string) ([]Elem, error) {
	return ReadWith(txt, DefaultParser)
}

func ReadWith(txt string, p Parser) ([]Elem, error) {
	var segments []Elem
	ctx := &ParseCtx{
		Buf: []byte(txt),
	}
	for {
		e, err := p.Parse(ctx)

		if err != nil {
			if err == io.EOF {
				break
			}
			return segments, err
		}

		segments = append(segments, e)
	}
	return segments, nil
}
