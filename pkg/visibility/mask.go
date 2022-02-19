package visibility

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

type MaskView struct {
	Store    zettel.Storage
	Tolerate []string
}

func (v MaskView) Mask(z zettel.Zettel) (zettel.Zettel, error) {
	return z.Rebuild(func(b zettel.Builder) error {
		var str strings.Builder
		elems, err := scan.Elements(v.Store, z.Readme().Text)
		if err != nil {
			return err
		}

		for _, el := range elems {
			txt := el.Span.String()

			if el.Type != scan.ItemRefbox {
				str.WriteString(txt)
				continue
			}

			scn := bufio.NewScanner(strings.NewReader(txt))
			scn.Scan()
			str.WriteString(fmt.Sprintln(scn.Text()))
			for scn.Scan() {
				line := scn.Text()
				zets, err := v.Store.Resolve(line)
				if err != nil {
					return err
				}
				if len(zets) != 1 {
					return fmt.Errorf("ambiguous zet ref not allowed in refbox: %q; %#v", line, zets)
				}

				z := zets[0]
				viz := Visible(z, v.Tolerate)
				if !viz {
					log.Printf("masking %s", z)
					z = Masked(z)
				}
				str.WriteString(zettel.MustFmt(z, "* {{.Id}}  {{.Title}}\n"))
			}
		}
		b.Text(str.String())
		return nil
	})
}

func Masked(z zettel.Zettel) zettel.Zettel {
	z2, e := z.Rebuild(func(b zettel.Builder) error {
		b.Title("MASKED")
		b.Metadata().Labels = make(zettel.Labels)
		b.Metadata().Link = nil
		b.Id(strings.Repeat("I", len(z.Id())))
		return nil
	})
	if e != nil {
		return z
	}
	return z2
}

func spanText(span scan.BlkSpan) string {
	lines := lines(span.Input)
	if span.Start == -1 {
		return span.Input
	}
	return strings.Join(lines[span.Start:span.Pos], "\n")
}

func lines(s string) []string {
	lins := make([]string, len(s)/80)
	scn := bufio.NewScanner(strings.NewReader(s))
	for scn.Scan() {
		lins = append(lins, scn.Text())
	}
	return lins
}
