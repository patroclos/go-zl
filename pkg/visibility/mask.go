package visibility

import (
	"strings"

	"git.jensch.dev/joshua/zl/pkg/zettel"
	"git.jensch.dev/joshua/zl/pkg/zettel/elemz"
)

type MaskView struct {
	Store    zettel.Storage
	Tolerate []string
}

func (v MaskView) Mask(z zettel.Z) (zettel.Z, error) {
	return z.Rebuild(func(b zettel.Builder) error {
		var str strings.Builder

		txt := z.Readme().Text
		boxes := elemz.Refboxes(txt)

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

			span := box.Span()
			if pos < span.Start {
				str.WriteString(txt[pos:span.Start])
			}
			str.WriteString(box.String())
			pos = span.End

		}
		if pos < len(txt) {
			str.WriteString(txt[pos:])
		}
		b.Text(str.String())
		return nil
	})
}
