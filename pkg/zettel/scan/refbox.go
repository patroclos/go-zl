package scan

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

const (
	refLine = `\* (?:<.+>$|[0-9]{6}-[a-zA-Z0-9]{4}  .+$)`
)

var regex = regexp.MustCompile(fmt.Sprintf(`(?m)^([a-zA-Z\[][^:\n]*):\n((?:%s)(?:\n%s)*)(\n\+ .*$(?:\n  .*$)*)?`, refLine, refLine))

type Refbox struct {
	Rel        string
	Refs       []string
	Extra      []string
	Start, End int
}

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

func All(txt string) []Refbox {
	matches := regex.FindAllStringSubmatchIndex(txt, -1)

	boxes := make([]Refbox, len(matches))

	for i, match := range matches {
		b := &boxes[i]
		b.Start, b.End = match[0], match[1]
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
