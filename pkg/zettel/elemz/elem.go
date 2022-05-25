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
	var elements []Elem
	ctx := &ParseCtx{Buf: []byte(txt)}

	for ctx.Pos < len(ctx.Buf) {
		e, err := p.Parse(ctx)

		if err == io.EOF {
			break
		}

		if err != nil {
			return elements, err
		}

		elements = append(elements, e)
	}
	return elements, nil
}
