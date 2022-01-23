package scan

import (
	"regexp"
	"strings"

	"jensch.works/zl/pkg/zettel"
)

func Refs(text string) []string {
	reg := regexp.MustCompile(`\[.+\]\((.+)\)`)
	matches := reg.FindAllStringSubmatch(text, -1)
	results := make([]string, 0, 8)
	for _, m := range matches {
		id := strings.Trim(m[1], " /")
		results = append(results, id)
	}
	reg = regexp.MustCompile(`\* ([a-zA-Z0-9-]+)  .*`)
	matches = reg.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		id := m[1]
		results = append(results, id)
	}

	return results
}

type RefZet struct {
	zettel.Zettel
	txt    string
	target zettel.Zettel
	link   zettel.Zettel
}

// link may be nil, target should always be valid
func (rz RefZet) Ref() (txt string, target zettel.Zettel, link zettel.Zettel) {
	return rz.txt, rz.target, rz.link
}
