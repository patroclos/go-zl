package visibility_test

import (
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/visibility"
	"jensch.works/zl/pkg/zettel"
)

func TestMaskTolerates(t *testing.T) {
	st, err := storage.NewStore(memfs.New())
	if err != nil {
		t.Fatal(err)
	}
	a, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id("aaa")
		b.Title("Topic")
		b.Metadata().Labels["zl/taint"] = "hidden"
		return nil
	})
	b, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id("bbb")
		b.Title("List")
		b.Text("Refs:\n* aaa  Topic")
		return nil
	})
	st.Put(a)
	st.Put(b)

	for _, tolerate := range []string{"hidden", "*", "xyz,hidden"} {
		split := strings.Split(tolerate, ",")

		mask, err := visibility.MaskView{
			Store:    st,
			Tolerate: split,
		}.Mask(b)
		if err != nil {
			t.Fatal(err)
		}

		if mask.Readme().Text != b.Readme().Text {
			t.Errorf("expected %q, got %q", b.Readme().Text, mask.Readme().Text)
		}
	}
}

func TestMaskRejects(t *testing.T) {
	st, err := storage.NewStore(memfs.New())
	if err != nil {
		t.Fatal(err)
	}
	a, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id("aaa")
		b.Title("Topic")
		b.Metadata().Labels["zl/taint"] = "hidden"
		return nil
	})
	b, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id("bbb")
		b.Title("List")
		b.Text("Refs:\n* aaa  Topic")
		return nil
	})
	st.Put(a)
	st.Put(b)

	for _, tolerate := range []string{"", "xyz"} {
		split := strings.Split(tolerate, ",")

		mask, err := visibility.MaskView{
			Store:    st,
			Tolerate: split,
		}.Mask(b)
		if err != nil {
			t.Fatal(err)
		}

		expect := "Refs:\n* III  MASKED"
		if mask.Readme().Text != expect {
			t.Errorf("expected %q, got %q", expect, mask.Readme().Text)
		}
	}
}

func TestMaskHalv(t *testing.T) {
	st, err := storage.NewStore(memfs.New())
	if err != nil {
		t.Fatal(err)
	}
	a, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id("aaa")
		b.Title("Topic")
		b.Metadata().Labels["zl/taint"] = "hidden"
		return nil
	})
	b, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id("bbb")
		b.Title("List")
		b.Text("Refs:\n* aaa  Topic\n* ccc  Topic")
		return nil
	})
	c, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id("ccc")
		b.Title("Topic")
		return nil
	})
	st.Put(a)
	st.Put(b)
	st.Put(c)

	for _, tolerate := range []string{"", "xyz"} {
		split := strings.Split(tolerate, ",")

		mask, err := visibility.MaskView{
			Store:    st,
			Tolerate: split,
		}.Mask(b)
		if err != nil {
			t.Fatal(err)
		}

		expect := "Refs:\n* III  MASKED\n* ccc  Topic"
		if mask.Readme().Text != expect {
			t.Errorf("expected %q, got %q", expect, mask.Readme().Text)
		}
	}
}
