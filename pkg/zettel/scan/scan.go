package scan

import (
	"bufio"
	"io"
	"strings"

	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

type Scanner interface {
	Scan(r io.Reader) <-chan zettel.Zettel
}

type listScanner struct {
	z storage.Zetteler
}

func ListScanner(z storage.Zetteler) Scanner {
	return listScanner{
		z: z,
	}
}

func (p listScanner) Scan(r io.Reader) <-chan zettel.Zettel {
	c := make(chan zettel.Zettel)
	go scan(c, p.z, r)
	return c
}

func scan(c chan<- zettel.Zettel, st storage.Zetteler, r io.Reader) {
	defer close(c)
	scn := bufio.NewScanner(r)
	for scn.Scan() {
		line := scn.Text()
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimLeft(line, " \t")

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		id := fields[0]

		zl, err := st.Zettel(id)
		if err != nil {
			continue
		}
		c <- zl
	}
}
