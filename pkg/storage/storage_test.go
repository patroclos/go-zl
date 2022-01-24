package storage

import (
	"fmt"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"jensch.works/zl/pkg/zettel"
)

func TestGitInit(t *testing.T) {
	var wt, dot billy.Filesystem
	wt = memfs.New()
	dot, _ = wt.Chroot(".git")

	s := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	_, err := git.Init(s, wt)

	if err != nil {
		t.Errorf("failed to make repo: %v", err)
	}
}

func TestNewStore(t *testing.T) {
	dir := memfs.New()
	st, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	if st.dir != dir {
		t.Error("NewStore(dir).dir is original dir")
	}

	if st.git == nil {
		t.Error("git repository is nil")
	}
}

func TestStore_Put(t *testing.T) {
	dir := memfs.New()
	st, _ := NewStore(dir)

	zl, _ := zettel.Construct(testZetConstructor)

	if err := st.Put(zl); err != nil {
		t.Fatal(err)
	}

	stat, err := dir.Lstat(fmt.Sprintf("%s/README.md", zl.Id()))
	if err != nil {
		t.Fatal(err)
	}
	if int(stat.Size()) != len(testContent) {
		t.Errorf("unexpected README.md size %d, expected %d", stat.Size(), len(testContent))
	}

	log, _ := st.git.Log(&git.LogOptions{All: true})
	commit, err := log.Next()
	if err != nil {
		t.Fatal("no git commits")
	}
	_ = commit

}

const testTitle = "Hello, Grid!"
const testContent = `Bla bla`

func testZetConstructor(b zettel.Builder) error {
	b.Title(testTitle)
	b.Text(testContent)
	return nil
}
