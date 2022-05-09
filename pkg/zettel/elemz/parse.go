package elemz

import (
	"fmt"
	"io"
	"log"
)

var DefaultParser *OneOfParser

func init() {
	DefaultParser = &OneOfParser{
		Parsers: []Parser{
			&refboxParser{},
			&codeParser{fence: "```"},
			&textParser{},
		},
	}
}

type Parser interface {
	Parse(*ParseCtx) (e Elem, err error)
}

type ParseCtx struct {
	Buf []byte
	Pos int
}

type OneOfParser struct {
	Parsers []Parser
}

func NewCompoundParser(parsers ...Parser) *OneOfParser {
	return &OneOfParser{parsers}
}

func (c *OneOfParser) Parse(ctx *ParseCtx) (e Elem, err error) {
	if ctx.Pos >= len(ctx.Buf) {
		return nil, io.EOF
	}
	for _, p := range c.Parsers {
		e, err := p.Parse(ctx)
		if err != nil {
			continue
		}

		if e == nil {
			log.Panicf("parser returned `nil,nil`. this is a bug, go fix it %T", p)
		}

		return e, nil
	}
	return nil, fmt.Errorf("none of the parsers matched %v ctx:%q", c.Parsers, ctx.Buf[ctx.Pos])
}
