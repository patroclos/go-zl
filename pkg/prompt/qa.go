package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"jensch.works/zl/pkg/zettel/elemz"
)

type QAPrompt struct {
	Q    string
	A    string
	span elemz.Span
}

const QAType elemz.ElemType = "zl/prompt.qa"

func (p *QAPrompt) Span() elemz.Span { return p.span }

func (p *QAPrompt) ElemType() elemz.ElemType { return QAType }

func (p QAPrompt) String() string {
	return fmt.Sprintf("Q.  %s\nA.  %s", p.Q, p.A)
}

type parseQA struct{}

func (p *parseQA) Parse(ctx *elemz.ParseCtx) (e elemz.Elem, err error) {
	scn := bufio.NewScanner(bytes.NewReader(ctx.Buf[ctx.Pos:]))
	if !scn.Scan() || !strings.HasPrefix(scn.Text(), "Q. ") {
		return nil, fmt.Errorf("'Q. ' prefix not found")
	}

	q := scn.Text()[3:]
	qlen := len(scn.Text())
	if len(q) == 0 {
		return nil, fmt.Errorf("empty question")
	}

	if !scn.Scan() {
		return nil, fmt.Errorf("parseQA: no line after q: %w", scn.Err())
	}

	span := elemz.Span{
		Start: ctx.Pos,
		End:   ctx.Pos + qlen + 1,
	}

	if !strings.HasPrefix(scn.Text(), "A. ") {
		return &QAPrompt{
			Q:    q,
			span: span,
		}, nil
	}
	span.End += len(scn.Bytes())

	ctx.Pos = span.End
	return &QAPrompt{
		Q:    q,
		A:    scn.Text()[3:],
		span: span,
	}, nil
}
