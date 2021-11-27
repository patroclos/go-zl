package scan_test

import (
	"fmt"
	"strings"
	"testing"

	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/storage/memory"
	"jensch.works/zl/pkg/zettel/scan"
)

func TestWithEmpty(t *testing.T) {
	txt := "20211111-111111-test"

	scn := scan.ListScanner(storage.Empty)
	ch := scn.Scan(strings.NewReader(txt))
	_, ok := <-ch

	if ok {
		t.Fatal("invalid scan")
	}
}

func TestMemoryOne(t *testing.T) {
	st := memory.NewStorage()
	zl := st.NewZettel("blabla")

	id := zl.Id()
	txt := fmt.Sprintf("* %s  blabla", id)

	scn := scan.ListScanner(st)
	ch := scn.Scan(strings.NewReader(txt))

	entry, ok := <-ch
	if !ok {
		t.Fatal("scan should have found blabla, found nothing")
	}

	if entry.Id() != id {
		t.Errorf("wanted to scan {Id => %s}, found %#v", id, entry)
	}

	entry, ok = <-ch
	if ok {
		t.Error("scanned too many", entry)
	}
}
