package zettel_test

import (
	"errors"
	"testing"

	"jensch.works/zl/pkg/zettel"
)

func TestParseLabelspec_default(t *testing.T) {
	txt := "zl/inbox=default-"
	spec, err := zettel.ParseLabelspec(txt)
	if err != nil {
		t.Fatalf("Failed parsing %s to labelspec: %+v", txt, err)
	}

	if spec.MatchLabel != "zl/inbox" {
		t.Error("invalid MatchLabel value", spec.MatchLabel)
	}

	if spec.MatchValue != "default" {
		t.Error("invalid MatchValue ", spec.MatchValue)
	}

	if !spec.Negated {
		t.Error("invalid Negated value", spec.Negated)
	}
}

func TestParseLabelspec_InvalidFormatErrors(t *testing.T) {
	cases := []string{
		"",
		"-=",
		"--",
		"/",
		"//",
		"//a",
		"a//b",
		"ab=",
	}

	for _, txt := range cases {
		spec, err := zettel.ParseLabelspec(txt)
		if err == nil {
			t.Errorf("Excpected %s to generate ErrInvalidSpecFormat. spec was %#v", txt, spec)
		}

		if !errors.Is(err, zettel.ErrInvalidSpecFormat) {
			t.Errorf("Invalid error for %s. Expected ErrInvalidSpecFormat got %v", txt, err)
		}
	}
}

func TestParseLabelspec(t *testing.T) {
	cases := []struct {
		txt string
		err error
		res *zettel.Labelspec
	}{
		{"zl/inbox=hot", nil, &zettel.Labelspec{"zl/inbox", "hot", false}},
		{"asdf", nil, &zettel.Labelspec{"asdf", "", false}},
		{"bsdf-", nil, &zettel.Labelspec{"bsdf", "", true}},
		{" oops", zettel.ErrInvalidSpecFormat, nil},
		{"aga!in", zettel.ErrInvalidSpecFormat, nil},
		{"aga!in", zettel.ErrInvalidSpecFormat, nil},
		{"/blub", zettel.ErrInvalidSpecFormat, nil},
		{"test/er/iNo0=!yeah", nil, &zettel.Labelspec{"test/er/iNo0", "!yeah", false}},
	}

	for _, cas := range cases {
		spec, err := zettel.ParseLabelspec(cas.txt)
		if cas.err == nil && err != nil {
			t.Errorf("Expected to parse %s but got %w", cas.txt, err)
			continue
		}

		if cas.err != nil {
			if !errors.Is(err, cas.err) {
				t.Errorf("Expected error of type %v, but got %w", cas.err, err)
			}
			continue
		}

		if cas.res == nil {
			continue
		}

		if cas.res.MatchLabel != spec.MatchLabel || cas.res.MatchValue != spec.MatchValue || cas.res.Negated != spec.Negated {
			t.Errorf("Expected %+v and %+v to match equal", cas.res, spec)
		}
	}
}

func TestLabelSpec_Met(t *testing.T) {
	spec := zettel.Labelspec{"zl/inbox", "default", false}
	if !spec.Match(zettel.Labels(map[string]string{"zl/inbox": "default"})) {
		t.Fatal("zl/inbox=default not met by zl/inbox: default")
	}
}

func TestInspiredLabelspec(t *testing.T) {
	spec, err := zettel.ParseLabelspec("zl/inbox=default-")
	if err != nil {
		t.Fatal(err)
	}
	if spec.MatchLabel != "zl/inbox" {
		t.Error("unexpected label", spec.MatchLabel)
	}
	if spec.MatchValue != "default" {
		t.Error("unexpected value", spec.MatchValue)
	}
	if !spec.Negated {
		t.Error("unexpected unnegated")
	}
}
