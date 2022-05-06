package elemz

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Elem interface {
	Spaner
	fmt.Stringer
}

type Spaner interface {
	Span() Span
}

type Span struct {
	Start int
	End   int
}

func Read(txt string) ([]Elem, error) {
	var res []Elem
	for seg := range ReadChan(strings.NewReader(txt)) {
		if seg.Err != nil {
			return res, seg.Err
		}
		res = append(res, seg.Elem)
	}
	return res, nil
}

type Segment struct {
	Elem Elem
	Err  error
}

func ReadChan(r io.Reader) <-chan Segment {
	c := make(chan Segment)
	go func() {
		defer close(c)
		buf := make([]byte, 0, 4096)
		pos := 0
		for {
			rbuffer := make([]byte, 4096)
			l, err := r.Read(rbuffer)
			if err == io.EOF {
				break
			}

			buf = append(buf, rbuffer[:l]...)

			e, adv, err := DefaultParser.Parse(ParseCtx{
				Buf: buf,
				Pos: pos,
			})
			pos += adv

			if e != nil {
				c <- Segment{Elem: e}
				log.Println(e, pos)
			}
			if err == io.EOF {
				break
			}

			if err != nil {
				c <- Segment{Err: err}
				return
			}
		}

		if pos < len(buf) {
			c <- Segment{
				Elem: &Text{
					Content: string(buf[pos:]),
					span:    Span{pos, pos + pos},
				},
			}
		}
	}()
	return c
}
