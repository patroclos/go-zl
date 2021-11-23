package filesystem

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	"jensch.works/zl/pkg/storage/memory"
	z "jensch.works/zl/pkg/zettel"
)

type Storage struct {
	ZettelStorage
}

type ZettelStorage struct {
	Directory string
}

type Zettel struct {
	memory.Zettel
	s     *ZettelStorage
	exist bool
}

func (zs *ZettelStorage) HasZettel(id z.Id) bool {
	_, err := zs.Zettel(id)
	return err == nil
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
	pz := path.Join(zl.Directory, string(z.Id()))
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
	readmePath := path.Join(zl.Directory, string(id), "README.md")
	metaPath := path.Join(zl.Directory, string(id), "meta.yaml")

	f, err := os.Open(readmePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scn := bufio.NewScanner(f)
	if !scn.Scan() {
		if err := scn.Err(); err != nil {
			return nil, err
		}
		return nil, errNoTitle
	}

	title := strings.TrimPrefix(scn.Text(), "# ")

	_, err = f.Seek(int64(len(title)+2), 0)
	if err != nil {
		return nil, err
	}

	rest, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var meta z.MetaInfo
	f, err = os.Open(metaPath)

	if err == nil {
		bytes, err := ioutil.ReadAll(f)
		if err == nil {
			err = yaml.Unmarshal(bytes, &meta)
			if err != nil {
				log.Printf("Error reading meta of %s:\n%v", title, err)
			}
		}
	}

	model := memory.CreateZettel(id, title, strings.TrimLeft(string(rest), "\n"), time.Now())
	metap, err := model.Metadata()
	if err != nil {
		panic(err)
	}
	*metap = meta

	zettel := Zettel{
		Zettel: model,
		s:      zl,
		exist:  true,
	}

	return &zettel, nil
}


func (zs ZettelStorage) ForEach(fn func(z z.Zettel) error) error {
	files, err := ioutil.ReadDir(zs.Directory)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		zettel, err := zs.Zettel(z.Id(f.Name()))
		if err != nil {
			continue
		}
		if err = fn(zettel); err != nil {
			return err
		}
	}

	return nil
}
