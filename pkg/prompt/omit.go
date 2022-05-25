package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	"jensch.works/zl/pkg/zettel/elemz"
)

type OmitPrompt struct {
	Fragments []string
	Holes     []string
	span      elemz.Span
}

const OmitType elemz.ElemType = "zl/prompt.omit"

func (p *OmitPrompt) Span() elemz.Span {
	return p.span
}

func (p *OmitPrompt) ElemType() elemz.ElemType { return OmitType }

func (p *OmitPrompt) String() string {
	sb := strings.Builder{}
	for i, f := range p.Fragments {
		sb.WriteString(f)

		if len(p.Holes) > i {
			sb.WriteString("{{")
			sb.WriteString(p.Holes[i])
			sb.WriteString("}}")
		}
	}
	return sb.String()
}

type parseOmit struct{}

func (_ *parseOmit) Parse(ctx *elemz.ParseCtx) (elemz.Elem, error) {
	scn := bufio.NewScanner(bytes.NewReader(ctx.Buf[ctx.Pos:]))

	if !scn.Scan() || !strings.Contains(scn.Text(), "{{") || !strings.Contains(scn.Text(), "}}") {
		return nil, fmt.Errorf("invalid begin: %q", scn.Text())
	}

	frags, holes, err := readOmit(scn.Text())
	if err != nil {
		return nil, err
	}

	p := &OmitPrompt{
		Fragments: frags,
		Holes:     holes,
	}
	log.Printf("omit: <%v, %v> %q\n", frags, holes, p.String())

	start := ctx.Pos
	ctx.Pos += len(p.String())
	return &OmitPrompt{
		Fragments: frags,
		Holes:     holes,
		span: elemz.Span{
			Start: start,
			End:   ctx.Pos,
		},
	}, err
}

func readOmit(txt string) (fragments []string, holes []string, err error) {
	open := 0
	clos := 0

	var hol strings.Builder
	var frag strings.Builder
	for _, char := range txt {
		switch open {
		case 0:
			if char == '{' {
				open = 1
				continue
			}
			frag.WriteRune(char)
			continue
		case 1:
			if char == '{' {
				open = 2
				continue
			}
			frag.WriteRune('{')
			frag.WriteRune(char)
			open = 0
			continue
		case 2:
			// inside hole, waiting for '}}'
			switch clos {
			case 0:
				switch char {
				case '}':
					clos = 1
				default:
					hol.WriteRune(char)
				}
				continue
			case 1:
				switch char {
				case '}':
					fragments = append(fragments, frag.String())
					holes = append(holes, hol.String())
					frag.Reset()
					hol.Reset()
					clos = 0
					open = 0
				default:
					hol.WriteRune('}')
					hol.WriteRune(char)
					clos = 0
				}
				continue
			}
		}
	}

	if frag.Len() > 0 {
		fragments = append(fragments, frag.String())
	}
	log.Println("last frag", txt)

	return fragments, holes, nil
}
