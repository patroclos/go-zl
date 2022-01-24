package zettel

import "testing"

func TestZetType(t *testing.T) {
	zl := &zet{}
	var _ Zettel = zl
}

func TestConstructEmpty(t *testing.T) {
	got, err := Construct(func(_ Builder) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal(got, err)
	}

	if got.Id() == "" {
		t.Error("expecting non-empty result from Id()")
	}
}
