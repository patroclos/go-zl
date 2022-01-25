package storage

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"gopkg.in/yaml.v2"
	"jensch.works/zl/pkg/zettel"
)

func NewStore(dir billy.Filesystem) (zettel.Storage, error) {
	return newStore(dir)
}

func newStore(dir billy.Filesystem) (*zetStore, error) {
	dotGit, _ := dir.Chroot(".git")
	gitStorage := filesystem.NewStorage(dotGit, cache.NewObjectLRUDefault())

	repo, err := git.Open(gitStorage, dir)
	if err != nil {
		repo, err = git.Init(gitStorage, dir)
	}

	if err != nil {
		return nil, err
	}

	st := &zetStore{
		dir: dir,
		git: repo,
		rw:  new(sync.RWMutex),
	}

	return st, nil
}

type zetStore struct {
	dir billy.Filesystem
	git *git.Repository
	rw  *sync.RWMutex
}

func (zs *zetStore) Zettel(id string) (zettel.Zettel, error) {
	chr, err := zs.dir.Chroot(id)
	if err != nil {
		return nil, err
	}

	return zettel.Read(id, chr)
}

func (zs *zetStore) Resolve(query string) (zettel.Zettel, error) {
	if zl, err := zs.Zettel(query); err == nil {
		return zl, nil
	}

	// TODO: match on filename (eg. resolving a full readme path to the zettel)
	// TODO: match on strings.Contains(Title(), query)
	// TODO: match on special queries like @last

	return nil, fmt.Errorf("couldn't resolve %s", query)
}

func (zs *zetStore) Put(zl zettel.Zettel) error {
	id := zl.Id()

	zs.rw.Lock()
	defer zs.rw.Unlock()

	tree, err := zs.git.Worktree()
	if err != nil {
		return err
	}

	status, err := tree.Status()
	if err != nil {
		return err
	}

	if !status.IsClean() {
		return fmt.Errorf("git worktree unclean")
	}

	zs.dir.Remove(id)

	if err := zs.dir.MkdirAll(id, 0); err != nil {
		return fmt.Errorf("can't create zettel dir: %w", err)
	}

	if err := writeReadme(zs, zl); err != nil {
		return err
	}

	if err := writeMeta(zs, id, zl.Metadata()); err != nil {
		return err
	}

	// create commit
	if err := tree.AddGlob(id); err != nil {
		return err
	}

	_, err = tree.Commit(fmt.Sprintf("put %s  %s", id, zl.Title()), &git.CommitOptions{})
	if err != nil {
		return err
	}

	return nil
}

func writeReadme(zs *zetStore, zl zettel.Zettel) error {
	id, r := zl.Id(), zl.Reader()
	if r == nil {
		return fmt.Errorf("reader is nil")
	}

	path := zs.dir.Join(id, "README.md")
	readme, err := zs.dir.Create(path)
	if err != nil {
		return fmt.Errorf("failed creating %s: %w", path, err)
	}

	io.Copy(readme, strings.NewReader(fmt.Sprintf("# %s\n\n", zl.Title())))
	defer readme.Close()
	_, err = io.Copy(readme, r)
	if err != nil {
		return fmt.Errorf("failed writing %s: %w", path, err)
	}
	return nil
}

func writeMeta(zs *zetStore, id string, mi *zettel.MetaInfo) error {
	path := zs.dir.Join(id, "meta.yaml")
	meta, err := zs.dir.Create(path)
	if err != nil {
		return fmt.Errorf("failed creating %s: %w", path, err)
	}
	defer meta.Close()

	mb, err := yaml.Marshal(mi)
	if err != nil {
		return fmt.Errorf("failed marshaling MetaInfo: %w", err)
	}

	_, err = io.Copy(meta, bytes.NewReader(mb))
	if err != nil {
		return fmt.Errorf("failed writing %s: %w", path, err)
	}
	return nil
}
