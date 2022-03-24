package storage

import (
	"io/fs"

	"github.com/go-git/go-billy/v5"
	"jensch.works/zl/pkg/zettel"
)

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
