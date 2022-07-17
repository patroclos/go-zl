package storage

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"jensch.works/zl/pkg/zettel"
)

func TestStoreType(t *testing.T) {
	dir := memfs.New()
	x, _ := NewStore(dir)
	var _ zettel.Storage = x
}

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
	st, err := newStore(dir)
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

func TestStore_PutNew(t *testing.T) {
	st, _ := newStore(memfs.New())
	zl, _ := zettel.Build(testZetConstructor)

	if err := st.Put(zl); err != nil {
		t.Fatal(err)
	}

	readme, err := st.dir.Open(fmt.Sprintf("%s/README.md", zl.Id()))
	defer readme.Close()
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := io.ReadAll(readme)
	got := string(bytes)

	expect := fmt.Sprintf("# %s\n\n%s", testTitle, testContent)

	if got != expect {
		t.Errorf("get '%s', expected '%s'", got, expect)
	}

	log, _ := st.git.Log(&git.LogOptions{All: true})
	commits := make([]*object.Commit, 0, 8)
	for {
		commit, err := log.Next()
		if err != nil {
			break
		}
		commits = append(commits, commit)
	}

	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}

	if !strings.Contains(commits[0].Message, zl.Id()) {
		t.Errorf("commit message didn't contain id: %s", commits[0].Message)
	}

	if !strings.Contains(commits[0].Message, zl.Readme().Title) {
		t.Errorf("commit message didn't contain title: %s", commits[0].Message)
	}
}

func TestStore_PutUpdate(t *testing.T) {
	st, _ := newStore(memfs.New())

	zl, _ := zettel.Build(testZetConstructor)
	altered, _ := zl.Rebuild(func(b zettel.Builder) error {
		b.Text("all new text")
		return nil
	})

	st.Put(zl)
	st.Put(altered)

	log, _ := st.git.Log(&git.LogOptions{All: true})
	commits := make([]*object.Commit, 0, 8)
	for {
		commit, err := log.Next()
		if err != nil {
			break
		}
		commits = append(commits, commit)
	}

	if len(commits) != 2 {
		t.Errorf("expected 2 commits, got %d", len(commits))
	}

	st2, _ := NewStore(st.dir)
	got, err := st2.Zettel(zl.Id())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil {
		t.Fatal("Zettel returned nil")
	}

	if got.Id() != zl.Id() {
		t.Fatal("id mismatch")
	}
}

// TODO: match on filename (eg. resolving a full readme path to the zettel)
// TODO: match on special queries like @last
func TestStoreResolve(t *testing.T) {
	st, _ := NewStore(memfs.New())

	zl, _ := zettel.Build(testZetConstructor)
	zl2, _ := zettel.Build(testZetConstructor)
	zl2, _ = zl2.Rebuild(func(b zettel.Builder) error {
		b.Title("Hello, foo!")
		return nil
	})

	st.Put(zl)
	st.Put(zl2)

	shouldMatch := []string{zl.Id(), zl.Id()[:4], zl2.Id(), "Grid", "foo"}

	for _, x := range shouldMatch {
		_, err := st.Resolve(x)
		if err != nil {
			t.Error(err)
		}
	}

	got, err := st.Resolve("Hello")
	if err != nil {
		t.Error(err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 matches, got %d", len(got))
	}
}

const testTitle = "Hello, Grid!"
const testContent = `Bla bla`

func testZetConstructor(b zettel.Builder) error {
	b.Title(testTitle)
	b.Text(testContent)
	return nil
}
