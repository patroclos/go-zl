package elemz

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

var DefaultParser Parser

func init() {
	DefaultParser = &CompoundParser{
		Parsers: []Parser{
			&codeParser{fence: "```"},
			&refboxParser{},
			&textParser{},
		},
	}
}

type Parser interface {
	// Parse can result in an error, miss or success.
	// When the input ends, io.EOF is returned as error, this may also be the
	// case when a result Elem is returned.
	Parse(ParseCtx) (e Elem, advance int, err error)
}

type ParseCtx struct {
	Buf []byte
	// bmarks, emarks?
	Pos int
}

type CompoundParser struct {
	Parsers []Parser
}

func NewCompoundParser(parsers ...Parser) *CompoundParser {
	return &CompoundParser{parsers}
}

func (c *CompoundParser) Parse(ctx ParseCtx) (e Elem, adv int, err error) {
	for _, p := range c.Parsers {
		e, a, err := p.Parse(ctx)
		adv += a
		if err != nil {
			return nil, adv, err
		}

		if e == nil {
			continue
		}
		return e, adv, err
	}
	return nil, adv, fmt.Errorf("none of the parsers matched %v", c.Parsers)
}

type textParser struct{}

func (p *textParser) Parse(ctx ParseCtx) (e Elem, adv int, err error) {
	scn := bufio.NewScanner(bytes.NewReader(ctx.Buf[ctx.Pos:]))

	var sb strings.Builder
	for scn.Scan() {
		tok := scn.Text()

		if strings.HasPrefix(tok, "```") || strings.HasSuffix(tok, ":") {
			return &Text{
				Content: sb.String(),
				span:    Span{ctx.Pos, ctx.Pos + sb.Len()},
			}, adv, nil
		}
		adv += len(scn.Text()) + 1
		sb.WriteString(fmt.Sprintln(tok))
	}
	if sb.Len() > 0 {
		return &Text{sb.String(), Span{ctx.Pos, ctx.Pos + sb.Len()}}, adv, nil
	}
	return nil, 0, io.EOF
}
