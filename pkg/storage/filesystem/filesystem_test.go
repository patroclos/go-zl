package filesystem_test

import (
	"testing"

	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/storage/filesystem"
)

func Test(t *testing.T) {
	store := filesystem.ZettelStorage{
		Directory: "testdata",
	}

	zs := storage.All(store)

	if len(zs) != 3 {
		t.Fail()
		t.Logf("Expected to find 3 zettel in testdata, found %d", len(zs))
	}
}
