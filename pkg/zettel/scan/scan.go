package scan

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"jensch.works/zl/pkg/zettel"
)

type Scanner interface {
	Scan(r io.Reader) <-chan zettel.Z
}

type listScanner struct {
	z Zettler
}

func ListScanner(z Zettler) Scanner {
	return listScanner{
		z: z,
	}
}

func (p listScanner) Scan(r io.Reader) <-chan zettel.Z {
	c := make(chan zettel.Z)
	go scan(c, p.z, r)
	return c
}

type Zettler interface {
	Zettel(id string) (zettel.Z, error)
}

func scan(c chan<- zettel.Z, st Zettler, r io.Reader) {
	defer close(c)
	scn := bufio.NewScanner(r)
	reg := regexp.MustCompile(`\[.+\]\((.+)\)`)
	for scn.Scan() {
		line := scn.Text()
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimLeft(line, " \t")

		matches := reg.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			id := strings.Trim(m[1], " /")
			if zl, err := st.Zettel(id); err == nil {
				c <- zl
			}
		}

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
