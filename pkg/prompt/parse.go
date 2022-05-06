package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"jensch.works/zl/pkg/zettel/elemz"
)

type parseQA struct{}

func (p *parseQA) Parse(ctx elemz.ParseCtx) (e elemz.Elem, adv int, err error) {
	scn := bufio.NewScanner(bytes.NewReader(ctx.Buf[ctx.Pos:]))
	scn.Scan()
	if !strings.HasPrefix(scn.Text(), "Q. ") {
		return nil, 0, nil
	}

	q := scn.Text()[3:]
	qlen := len(scn.Text())
	if len(q) == 0 {
		return nil, 0, nil
	}

	if !scn.Scan() {
		return nil, 0, fmt.Errorf("parseQA: no line after q: %w", scn.Err())
	}

	if !strings.HasPrefix(scn.Text(), "A. ") {
		return nil, 0, nil
	}

	adv = qlen + 1 + len(scn.Bytes())
	return &QAPrompt{Q: q, A: scn.Text()[3:], span: elemz.Span{Start: ctx.Pos, End: ctx.Pos + adv}}, adv, nil
}
