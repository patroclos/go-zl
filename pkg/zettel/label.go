package zettel

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidSpecFormat = errors.New("invalid spec format")
)

// Constraint for filtering knodes
type Labelspec struct {
	MatchLabel string
	MatchValue string
	Negated    bool
}

func RunSpecs(specs []Labelspec, labels Labels) bool {
	for _,spec := range specs {
		if !RunLabelspec(spec, labels) {
			return false
		}
	}
	return true
}

func RunLabelspec(ls Labelspec, labels Labels) bool {
	ignoreVal := ls.MatchValue == ""
	var val *string = nil
	for k, v := range labels {
		vv := v
		if k == ls.MatchLabel {
			val = &vv
		}
	}

	if ignoreVal && val != nil {
		return !ls.Negated
	}

	if !ignoreVal && val != nil && *val != ls.MatchValue {
		return ls.Negated
	}

	return val != nil
}

func ParseLabelspec(txt string) (*Labelspec, error) {
	if len(txt) == 0 {
		return nil, fmt.Errorf("spec is empty: %w", ErrInvalidSpecFormat)
	}

	negated := txt[0] == '-'
	if negated {
		txt = txt[1:]
	}

	if len(txt) == 0 {
		return &Labelspec{
			MatchLabel: "",
			MatchValue: "",
			Negated:    true,
		}, nil
	}

	isep := strings.Index(txt, "=")
	if isep == -1 {
		if err := validateLabelName(txt); err != nil {
			return nil, err
		}
		return &Labelspec{
			MatchLabel: txt,
			MatchValue: "",
			Negated:    negated,
		}, nil
	}

	comps := strings.SplitN(txt, "=", 2)
	if len(comps) != 2 {
		return nil, fmt.Errorf("len(comps) != 2: %w", ErrInvalidSpecFormat)
	}

	name, value := comps[0], comps[1]

	if err := validateLabelName(name); err != nil {
		return nil, err
	}

	if err := validateLabelValue(name, value); err != nil {
		return nil, err
	}

	return &Labelspec{
		MatchLabel: name,
		MatchValue: value,
		Negated:    negated,
	}, nil
}

func validateLabelName(comps string) error {
	regex := regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9\-\/]*)$`)
	if !regex.Match([]byte(comps)) {
		return fmt.Errorf("'%s' did not match the label regex %s: %w", comps, regex.String(), ErrInvalidSpecFormat)
	}

	for _, c := range []string{"/", "-"} {
		if strings.HasPrefix(comps, c) {
			return ErrInvalidSpecFormat
		}
	}

	if strings.HasSuffix(comps, "/") {
		return fmt.Errorf("label name has invalid suffing /: %w", ErrInvalidSpecFormat)
	}

	if strings.Contains(comps, "//") {
		prefix := "label name contains double slash: "

		idx := strings.Index(comps, "//")
		sb := strings.Builder{}
		sb.WriteString(comps)
		sb.WriteRune('\n')
		for i := 0; i < idx+len(prefix); i++ {
			sb.WriteRune(' ')
		}
		sb.WriteString("⬆️ ⬆️")

		return fmt.Errorf("%s%w\n%s", prefix, ErrInvalidSpecFormat, sb.String())
	}

	return nil
}

func validateLabelValue(name string, value string) error {
	if len(value) == 0 {
		return fmt.Errorf("%s value: %w", name, ErrInvalidSpecFormat)
	}

	return nil
}
