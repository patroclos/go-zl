package scan

import (
	"bufio"
	"io"
	"strings"

	"jensch.works/zl/pkg/zettel"
)

type Scanner interface {
	Scan(r io.Reader) <-chan zettel.Zettel
}

type listScanner struct {
	z Zettler
}

func ListScanner(z Zettler) Scanner {
	return listScanner{
		z: z,
	}
}

func (p listScanner) Scan(r io.Reader) <-chan zettel.Zettel {
	c := make(chan zettel.Zettel)
	go scan(c, p.z, r)
	return c
}

type Zettler interface {
	Zettel(id string) (zettel.Zettel, error)
}

func scan(c chan<- zettel.Zettel, st Zettler, r io.Reader) {
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
