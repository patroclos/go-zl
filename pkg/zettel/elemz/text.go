package elemz

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Text struct {
	Content string
	span    Span
}

const TextType ElemType = "zl/text"

func (e *Text) Span() Span         { return e.span }
func (e *Text) ElemType() ElemType { return TextType }
func (e *Text) String() string     { return e.Content }

type textParser struct{}

func (p *textParser) Parse(ctx *ParseCtx) (e Elem, err error) {
	scn := bufio.NewScanner(bytes.NewReader(ctx.Buf[ctx.Pos:]))

	var sb strings.Builder
	adv := 0
	for scn.Scan() {
		tok := scn.Text()

		if strings.HasPrefix(tok, "```") || strings.HasSuffix(tok, ":") {
			if sb.Len() == 0 {
				goto ok
			}
			start := ctx.Pos
			ctx.Pos += adv
			return &Text{
				Content: sb.String(),
				span:    Span{start, ctx.Pos},
			}, nil
		}

	ok:
		adv += len(tok) + 1
		sb.WriteString(fmt.Sprintln(tok))
	}
	if adv > 0 {
		start := ctx.Pos
		ctx.Pos += adv
		return &Text{sb.String(), Span{start, ctx.Pos}}, nil
	}
	return nil, io.EOF
}
