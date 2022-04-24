package visibility

import (
	"strings"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

type MaskView struct {
	Store    zettel.Storage
	Tolerate []string
}

func (v MaskView) Mask(z zettel.Z) (zettel.Z, error) {
	return z.Rebuild(func(b zettel.Builder) error {
		var str strings.Builder

		txt := z.Readme().Text
		boxes := scan.All(txt)

		pos := 0

		for _, box := range boxes {
			for i, ref := range box.Refs {
				zl, err := v.Store.Zettel(ref[:11])
				if err != nil {
					continue
				}
				viz := Visible(zl, v.Tolerate)
				if !viz {
					box.Refs[i] = "000000-MASK  MASKED"
				}
			}

			if pos < box.Start {
				str.WriteString(txt[pos:box.Start])
			}
			str.WriteString(box.String())
			pos = box.End

		}
		if pos < len(txt) {
			str.WriteString(txt[pos:])
		}
		b.Text(str.String())
		return nil
	})
}
