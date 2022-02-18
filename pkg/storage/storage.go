package storage

import (
	"fmt"
	"io/fs"
	"log"
	"strings"
	"sync"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"jensch.works/zl/pkg/storage/strutil"
	_ "jensch.works/zl/pkg/storage/strutil"
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

type iter struct {
	dir     billy.Filesystem
	files   []fs.FileInfo
	current zettel.Zettel
}

func (i *iter) Next() bool {
	if i.files == nil {
		files, err := i.dir.ReadDir("")
		if err != nil {
			return false
		}
		i.files = files
	}

	if len(i.files) == 0 {
		return false
	}

	var x fs.FileInfo = nil
	for len(i.files) > 0 {
		a, xs := i.files[0], i.files[1:]
		i.files = xs

		if a.IsDir() {
			x = a
			break
		}
	}

	if x == nil {
		return false
	}

	zroot, err := i.dir.Chroot(x.Name())
	if err != nil {
		// log.Println(err, x.Name())
		return i.Next()
	}
	zet, err := zettel.Read(x.Name(), zroot)
	if err != nil {
		// log.Println(err, x.Name())
		return i.Next()
	}

	i.current = zet

	return true
}

func (i *iter) Zet() zettel.Zettel {
	if i.current == nil {
		panic("you are cringe")
	}
	return i.current
}

func (zs *zetStore) Iter() zettel.Iterator {
	return &iter{dir: zs.dir}
}

func (zs *zetStore) Resolve(query string) ([]zettel.Zettel, error) {
	query = strings.TrimPrefix(query, "* ")
	log.Println(query)
	if zl, err := zs.Zettel(query); err == nil {
		return []zettel.Zettel{zl}, nil
	}

	partialMatches := make([]zettel.Zettel, 0, 8)
	infos, err := zs.dir.ReadDir("")
	if err != nil {
		return nil, err
	}

	for _, x := range infos {
		if !x.IsDir() {
			continue
		}
		id := x.Name()
		ch, _ := zs.dir.Chroot(id)
		zet, err := zettel.Read(id, ch)
		if err != nil {
			continue
		}

		match := strutil.ContainsFold(fmt.Sprint(zet), query)
		if match {
			partialMatches = append(partialMatches, zet)
		}
	}
	if len(partialMatches) != 0 {
		return partialMatches, nil
	}

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

	chroot, err := zs.dir.Chroot(id)
	if err := zettel.Write(zl, chroot); err != nil {
		return err
	}

	// create commit
	if err := tree.AddGlob(id); err != nil {
		return err
	}

	_, err = tree.Commit(fmt.Sprintf("%s  %s", id, zl.Readme().Title), &git.CommitOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (zs *zetStore) Remove(zet zettel.Zettel) error {
	id := zet.Id()
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

	if _, err := tree.Remove(id); err != nil {
		return err
	}

	// create commit
	_, err = tree.Commit(fmt.Sprintf("REMOVE %s  %s", id, zet.Readme().Title), &git.CommitOptions{})
	if err != nil {
		return err
	}

	return nil
}
