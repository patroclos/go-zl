package elemz

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

const CodeType ElemType = "zl/code"

type Code struct {
	BlockParam string
	Code       string
	span       Span
}

func (e *Code) Span() Span         { return e.span }
func (e *Code) ElemType() ElemType { return CodeType }
func (e *Code) String() string     { return fmt.Sprintf("```%s\n%s\n```", e.BlockParam, e.Code) }

type codeParser struct {
	fence string
}

func (p *codeParser) Parse(ctx *ParseCtx) (e Elem, err error) {
	scn := bufio.NewScanner(bytes.NewReader(ctx.Buf[ctx.Pos:]))

	if !scn.Scan() || !strings.HasPrefix(scn.Text(), p.fence) {
		return nil, fmt.Errorf("no fence found")
	}
	lines := []string{scn.Text()}
	for scn.Scan() {
		lines = append(lines, scn.Text())
		if scn.Text() == p.fence {
			break
		}
	}

	if len(lines) == 1 || lines[len(lines)-1] != p.fence {
		return nil, fmt.Errorf("code-block not closed: %v", lines)
	}
	txt := strings.Join(lines, "\n")
	start := ctx.Pos
	ctx.Pos += len(txt)
	return &Code{
		BlockParam: strings.TrimPrefix(lines[0], p.fence),
		Code:       strings.Join(lines[1:len(lines)-1], "\n"),
		span:       Span{start, ctx.Pos},
	}, nil
}
