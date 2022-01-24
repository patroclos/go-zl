package storage

import (
	"fmt"
	"io/ioutil"
	"strings"
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

	readme, err := dir.Open(fmt.Sprintf("%s/README.md", zl.Id()))
	defer readme.Close()
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := ioutil.ReadAll(readme)
	got := string(bytes)

	expect := fmt.Sprintf("# %s\n\n%s", testTitle, testContent)

	if got != expect {
		t.Errorf("get '%s', expected '%s'", got, expect)
	}

	log, _ := st.git.Log(&git.LogOptions{All: true})
	commit, err := log.Next()
	if err != nil {
		t.Fatal("no git commits")
	}
	_ = commit

	if !strings.Contains(commit.Message, zl.Id()) {
		t.Errorf("commit message didn't contain id: %s", commit.Message)
	}

	if !strings.Contains(commit.Message, zl.Title()) {
		t.Errorf("commit message didn't contain title: %s", commit.Message)
	}
}

const testTitle = "Hello, Grid!"
const testContent = `Bla bla`

func testZetConstructor(b zettel.Builder) error {
	b.Title(testTitle)
	b.Text(testContent)
	return nil
}
