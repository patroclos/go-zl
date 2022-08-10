package search

import (
	"strings"

	"git.jensch.dev/zl/pkg/zettel"
)

type Params struct {
	Plain  string
	Labels map[string]zettel.Labelspec
}

func Query(q string) (*Params, error) {
	var b strings.Builder
	sep := false
	labels := make([]zettel.Labelspec, 0)

	for _, txt := range strings.Fields(q) {
		if strings.HasPrefix(txt, "label:") {
			spec, err := zettel.ParseLabelspec(txt[len("label:"):])
			if err != nil {
				return nil, err
			}
			labels = append(labels, spec)
			continue
		}
		if sep {
			b.WriteRune(' ')
		} else {
			sep = true
		}
		b.WriteString(txt)
	}

	labelMap := make(map[string]zettel.Labelspec, len(labels))
	for _, spec := range labels {
		labelMap[spec.MatchLabel] = spec
	}

	return &Params{
		Plain:  b.String(),
		Labels: labelMap,
	}, nil
}
