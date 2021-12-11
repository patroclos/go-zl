package search_test

import (
	"testing"

	"jensch.works/zl/pkg/search"
	"jensch.works/zl/pkg/storage/memory"
)

var (
	_ = memory.NewStorage()
	_ = search.SearchPattern{}
)

func TestResolveId(t *testing.T) {
}
