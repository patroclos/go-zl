package scan_test

import (
	"strings"
	"testing"

	"jensch.works/zl/pkg/storage"
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
