package storage

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
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

func FromEnv() (zettel.Storage, error) {
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		return nil, fmt.Errorf("no environment variable ZLPATH found")
	}
	dir := osfs.New(zlpath)
	return NewStore(dir)
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

func (zs *zetStore) Zettel(id string) (zettel.Z, error) {
	chr, err := zs.dir.Chroot(id)
	if err != nil {
		return nil, err
	}

	return zettel.Read(id, chr)
}

func (zs *zetStore) Iter() zettel.Iterator {
	return &iter{dir: zs.dir}
}

func (zs *zetStore) Resolve(query string) ([]zettel.Z, error) {
	query = strings.TrimPrefix(query, "* ")
	if zl, err := zs.Zettel(query); err == nil {
		return []zettel.Z{zl}, nil
	}

	if split := strings.Split(query, "  "); len(split) > 1 {
		if zl, err := zs.Zettel(split[0]); err == nil {
			return []zettel.Z{zl}, nil
		}
	}

	partialMatches := make([]zettel.Z, 0, 8)
	infos, err := zs.dir.ReadDir("")
	if err != nil {
		return nil, err
	}

	for _, x := range infos {
		if !x.IsDir() || x.Name() == ".git" {
			continue
		}
		id := x.Name()
		ch, _ := zs.dir.Chroot(id)
		zet, err := zettel.Read(id, ch)
		if err != nil {
			log.Println(err)
			continue
		}

		match := strutil.ContainsFold(zettel.MustFmt(zet, zettel.ListingFormat), query)
		if match {
			partialMatches = append(partialMatches, zet)
		}
	}
	if len(partialMatches) != 0 {
		return partialMatches, nil
	}

	return nil, fmt.Errorf("couldn't resolve %s", query)
}

func (zs *zetStore) Put(zl zettel.Z) error {
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

	id := zl.Id()
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

func (zs *zetStore) Remove(zet zettel.Z) error {
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
