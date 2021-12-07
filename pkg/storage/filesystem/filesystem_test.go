package filesystem_test

import (
	"testing"
	"time"

	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/storage/filesystem"
	"jensch.works/zl/pkg/zettel"
)

func TestTimeFromId(t *testing.T) {
	T := time.Now()
	format := "20060102-030405"
	id := zettel.Id(T.Format(format))
	parsed, err := filesystem.ParseTimeFromId(id)

	if err != nil {
		t.Fatal(err)
	}

	if tf, pf := T.Format(format), parsed.Format(format); tf != pf {
		t.Errorf("parsed equivalent of %s (%s), expected %s", pf, parsed.Format(time.RFC3339), tf)
	}
}

func TestTimestampIdRead(t *testing.T) {
	store := filesystem.ZettelStorage{
		Directory: "testdata",
	}

	zl, err := store.Zettel("20210831-010203-suffix")
	if err != nil {
		t.Fatal(err)
	}
	meta, err := zl.Metadata()
	if err != nil {
		t.Fatal("Zettel contains no metadata", zl)
	}

	expected := time.Date(2021, 8, 31, 01, 02, 03, 0, time.Local)
	if !meta.CreateTime.Equal(expected) {
		t.Error("expected equal", meta.CreateTime, expected)
	}
}

func TestHasZettel(t *testing.T) {
	store := filesystem.ZettelStorage{
		Directory: "testdata",
	}

	for _, i := range []string{"id1", "id2", "id3"} {
		if !store.HasZettel(zettel.Id(i)) {
			t.Errorf("expected to find zettel with id %s", i)
		}
	}
	zs := storage.All(store)

	if len(zs) < 3 {
		t.Fail()
		t.Logf("Expected to find 3 zettel in testdata, found %d", len(zs))
	}
}

func TestLinkMetadata(t *testing.T) {
	store := filesystem.ZettelStorage{
		Directory: "testdata",
	}

	zl, err := store.Zettel("lnkid1id2")
	if err != nil {
		t.Fatal("Dind't find expected link zettel 'lnkid1id2' in testdata", err)
		return
	}

	meta, err := zl.Metadata()
	if err != nil {
		t.Fatal("No metadata parsed for link zettel")
		return
	}
	if meta.Link == nil {
		t.Fatal("No link metadata parsed for link zettel")
		return
	}

	if meta.Link.A != "id1" || meta.Link.B != "id2" {
		t.Errorf("expected {from: id1, to: id2}, got {from: %s, to: %s}", meta.Link.A, meta.Link.B)
	}

	if l := len(meta.Link.Ctx); l != 1 {
		t.Errorf("Expected 1 context zettel, found %d", l)
		return
	}

	if meta.Link.Ctx[0] != "id3" {
		t.Errorf("expected Link.Ctx to be [id3], got [%s]", meta.Link.Ctx[0])
	}
}
