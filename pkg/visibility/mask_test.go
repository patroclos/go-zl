package visibility_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"git.jensch.dev/zl/pkg/storage"
	"git.jensch.dev/zl/pkg/visibility"
	"git.jensch.dev/zl/pkg/zettel"
)

func TestMaskTolerates(t *testing.T) {
	st, err := storage.NewStore(memfs.New())
	if err != nil {
		t.Fatal(err)
	}
	id1, id2 := zettel.MakeId(), zettel.MakeId()
	a, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(id1)
		b.Title("Topic")
		b.Metadata().Labels["zl/taint"] = "hidden"
		return nil
	})
	b, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(id2)
		b.Title("List")
		b.Text(fmt.Sprintf("Refs:\n* %s  Topic", id1))
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
	id1, id2 := zettel.MakeId(), zettel.MakeId()
	a, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(id1)
		b.Title("Topic")
		b.Metadata().Labels["zl/taint"] = "hidden"
		return nil
	})
	b, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(id2)
		b.Title("List")
		b.Text(fmt.Sprintf("Refs:\n* %s  Topic", id1))
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

		expect := "Refs:\n* 000000-MASK  MASKED"
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
	ids := make([]string, 3)
	for i := range ids {
		ids[i] = zettel.MakeId()
	}
	a, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(ids[0])
		b.Title("Topic")
		b.Metadata().Labels["zl/taint"] = "hidden"
		return nil
	})
	b, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(ids[1])
		b.Title("List")
		b.Text(fmt.Sprintf("Refs:\n* %s  Topic\n* %s  Topic", ids[0], ids[2]))
		return nil
	})
	c, _ := zettel.Build(func(b zettel.Builder) error {
		b.Id(ids[2])
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

		expect := fmt.Sprintf("Refs:\n* 000000-MASK  MASKED\n* %s  Topic", ids[2])
		if mask.Readme().Text != expect {
			t.Errorf("expected %q, got %q", expect, mask.Readme().Text)
		}
	}
}
