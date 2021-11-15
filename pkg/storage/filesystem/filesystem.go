package filesystem

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/storage/memory"
	z "jensch.works/zl/pkg/zettel"
)

type Storage struct {
	ZettelStorage
}

type ZettelStorage struct {
	dir string
}

type Zettel struct {
	memory.Zettel
	s     *ZettelStorage
	exist bool
}

// return a filesystem storage based on the directory configured in XDG_USER_CONFIG
func NewStorage(dir string) (storage.Storer, error) {
	return &ZettelStorage{dir: dir}, nil
}

func (zs *ZettelStorage) NewZettel(title string) z.Zettel {
	zl := memory.NewZettel(title)
	zz := Zettel{
		Zettel: zl,
		s:      zs,
	}
	return &zz
}

func (zl *ZettelStorage) SetZettel(z z.Zettel) error {
	pz := path.Join(zl.dir, string(z.Id()))
	err := os.MkdirAll(pz, 0755)
	if err != nil {
		return err
	}

	txt, err := z.Text()
	if err != nil {
		return err
	}

	md := fmt.Sprintf("%s\n\n%s", z.Title(), txt)
	err = ioutil.WriteFile(path.Join(pz, "README.md"), []byte(md), 644)
	if err != nil {
		return err
	}

	return nil
}

var errNoTitle = errors.New("no title")

func (zl *ZettelStorage) Zettel(id z.Id) (z.Zettel, error) {
	readmePath := path.Join(zl.dir, string(id), "README.md")

	f, err := os.Open(readmePath)
	if err != nil {
		return nil, err
	}

	scn := bufio.NewScanner(f)
	if !scn.Scan() {
		if err := scn.Err(); err != nil {
			return nil, err
		}
		return nil, errNoTitle
	}

	title := scn.Text()
	_, err = f.Seek(int64(len(title) + 2), 0)
	if err != nil {
		return nil, err
	}

	rest, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	model := memory.CreateZettel(id, title, string(rest), time.Now())

	zettel := Zettel{
		Zettel: model,
		s: zl,
		exist: true,
	}

	return &zettel, nil
}

func (zl *ZettelStorage) IterZettel() (storage.ZettelIter, error) {
	return nil, errors.New("no")
}
