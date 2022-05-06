package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"jensch.works/zl/pkg/zettel/elemz"
)

// The prompt package contains everything related to Q.A. prompts and omissables
// the existance of prompts puts pressure on a view-over-edit mentality to experience
// properly embedded prompts.

type QAPrompt struct {
	Q    string
	A    string
	span elemz.Span
}

func (p *QAPrompt) Span() elemz.Span { return p.span }

func (p QAPrompt) String() string {
	return fmt.Sprintf("Q.  %s\nA.  %s", p.Q, p.A)
}

type OmitPrompt struct {
	Fragments []string
	Holes     []Hole
}

func (p OmitPrompt) String() string {
	sb := strings.Builder{}
	for i, f := range p.Fragments {
		sb.WriteString(f)

		if len(p.Holes) > i {
			sb.WriteString("{{")
			sb.WriteString(p.Holes[i].Text)
			sb.WriteString("}}")
		}
	}
	return sb.String()
}

type Hole struct {
	Text string
}

type Stringer interface {
	String() string
}

type EmbeddedPrompt interface {
	Stringer
}

func Parser() elemz.Parser {
	return elemz.NewCompoundParser(&parseQA{})
}

func ExtractAll(txt string) []EmbeddedPrompt {
	scn := bufio.NewScanner(bytes.NewReader([]byte(txt)))

	results := make([]EmbeddedPrompt, 0, 4)

	var seenQ bool
	var q string
	for scn.Scan() {
		line := scn.Text()
		if seenQ {
			if !strings.HasPrefix(line, "A. ") {
				seenQ = false
				continue
			}

			results = append(results, QAPrompt{
				Q: q,
				A: line[3:],
			})
			continue
		}

		if strings.HasPrefix(line, "Q. ") {
			seenQ = true
			q = line[3:]
			continue
		}

		omit := OmitPrompt{
			Fragments: make([]string, 0, 8),
			Holes:     make([]Hole, 0, 8),
		}

		p := 0
		for p < len(line) {

			iopen := strings.Index(line[p:], "{{")
			if iopen == -1 {
				omit.Fragments = append(omit.Fragments, line[p:])
				break
			}
			iopen += p
			iclose := strings.Index(line[iopen:], "}}")

			if iclose == -1 {
				break
			}

			iclose += iopen

			if iopen+2 == iclose {
				p = p + iclose
				continue
			}

			if iclose > -1 {
				txt := line[iopen+2 : iclose]
				omit.Fragments = append(omit.Fragments, line[p:iopen])
				omit.Holes = append(omit.Holes, Hole{Text: txt})
				p = iclose + 2
			}
		}

		if len(omit.Holes) > 0 {
			results = append(results, omit)
		}

	}
	return results
}
