package elemz

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const (
	refLine             = `\* (?:<.+>$|[0-9]{6}-[a-zA-Z0-9]{4}  .+$)`
	RefboxType ElemType = "zl/refbox"
)

var regex = regexp.MustCompile(fmt.Sprintf(`(?m)^([a-zA-Z\[][^:\n]*):\n((?:%s)(?:\n%s)*)(\n\+ .*$(?:\n  .*$)*)?`, refLine, refLine))

type refboxParser struct {
}

func (p *refboxParser) Parse(ctx *ParseCtx) (e Elem, err error) {
	scn := bufio.NewScanner(bytes.NewReader(ctx.Buf[ctx.Pos:]))
	if !scn.Scan() {
		return nil, fmt.Errorf("no rel line")
	}
	rel := scn.Text()
	if !strings.HasSuffix(rel, ":") {
		// emit text
		// el := &Text{Content: rel, span: Span{Start: ctx.Pos, End: ctx.Pos + len(rel) + 1}}
		// return el, el.span.End - el.span.Start, nil
		return nil, fmt.Errorf("expected rel ':' suffix")
	}

	refs := []string{}
	extra := []string{}
	adv := len(rel) + 1
	for scn.Scan() {
		adv += len(scn.Bytes()) + 1
		if len(scn.Text()) == 0 {
			break
		}

		switch scn.Text()[0] {
		case '*':
			refs = append(refs, scn.Text())
		case '+':
			extra = append(extra, scn.Text())
		default:
			break
		}
	}

	for i, ref := range refs {
		ref = strings.TrimPrefix(ref, "* ")
		refs[i] = ref
		if strings.HasPrefix(ref, "<") {
			continue
		}
		if strings.Index(ref, "  ") == -1 {
			return nil, fmt.Errorf("non-ref found in refbox")
		}
	}

	if len(refs) == 0 {
		return nil, fmt.Errorf("need at least one ref")
	}

	start := ctx.Pos
	ctx.Pos = ctx.Pos + adv
	return &Refbox{
		Rel:     rel,
		Refs:    refs,
		Extra:   extra,
		BoxSpan: Span{start, ctx.Pos},
	}, nil
}

type Refbox struct {
	Rel     string
	Refs    []string
	Extra   []string
	BoxSpan Span
}

func (e *Refbox) Span() Span {
	return e.BoxSpan
}

func (e *Refbox) ElemType() ElemType { return RefboxType }

func (r Refbox) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s:\n", r.Rel))
	if len(r.Refs) > 0 {
		b.WriteString(fmt.Sprintf("* %s", r.Refs[0]))
		for _, x := range r.Refs[1:] {
			b.WriteString(fmt.Sprintf("\n* %s", x))
		}
	}

	if len(r.Extra) > 0 {
		b.WriteString(fmt.Sprintf("\n+ %s", r.Extra[0]))

		for _, x := range r.Extra[1:] {
			b.WriteString(fmt.Sprintf("\n  %s", x))
		}
	}
	return b.String()
}

func Refboxes(txt string) []Refbox {
	matches := regex.FindAllStringSubmatchIndex(txt, -1)

	boxes := make([]Refbox, len(matches))

	for i, match := range matches {
		b := &boxes[i]
		b.BoxSpan.Start, b.BoxSpan.End = match[0], match[1]
		b.Rel = txt[match[2]:match[3]]

		b.Refs = make([]string, 0, 8)
		refs := txt[match[4]:match[5]]
		for scn := bufio.NewScanner(strings.NewReader(refs)); scn.Scan(); {
			b.Refs = append(b.Refs, scn.Text()[2:])
		}

		if match[6] < 0 {
			continue
		}

		b.Extra = make([]string, 0, 2)
		extra := txt[match[6]:match[7]]
		scn := bufio.NewScanner(strings.NewReader(extra))
		scn.Scan()
		for scn.Scan() {
			if len(scn.Text()) < 2 {
				panic(fmt.Sprintf("extra too short: %q", scn.Text()))
			}
			b.Extra = append(b.Extra, scn.Text()[2:])
		}

	}
	return boxes
}
