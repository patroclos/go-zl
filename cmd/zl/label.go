package main

import (
	"fmt"
	"os"

	"github.com/go-clix/cli"
	"git.jensch.dev/joshua/go-zl/pkg/zettel"
)

func makeCmdLabel(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{
		Use:  "label [ref] [spec...]",
		Args: cli.ArgsMin(1),
		Run: func(cmd *cli.Command, args []string) error {
			zets, err := st.Resolve(args[0])
			if err != nil {
				return err
			}
			zet, err := pickOne(zets)
			if err != nil {
				return err
			}

			if isTerminal(os.Stdin) {
				if len(args) == 1 {
					fmt.Println(zet.Metadata().Labels)
					return nil
				}

				specs := make([]zettel.Labelspec, len(args)-1)
				orig := zet
				for i := range specs {
					spec, err := zettel.ParseLabelspec(args[i+1])
					if err != nil {
						return err
					}
					specs[i] = spec
					zet, err = stampLabels(zet, spec)
					if err != nil {
						return err
					}
				}

				if orig.Metadata().Equal(zet.Metadata()) {
					fmt.Println(zet)
					return nil
				}
				if err := st.Put(zet); err != nil {
					return err
				}
				fmt.Println(zet)
			}
			return nil
		},
	}
	return cmd
}

func stampLabels(z zettel.Z, spec zettel.Labelspec) (zettel.Z, error) {
	return z.Rebuild(func(b zettel.Builder) error {
		m := b.Metadata()
		val, ok := m.Labels[spec.MatchLabel]
		if spec.Negated {
			if !ok {
				return nil
			}
			if spec.MatchValue != "" {
				if val == spec.MatchValue {
					delete(m.Labels, spec.MatchLabel)
				}
				return nil
			}
			delete(m.Labels, spec.MatchLabel)
			return nil
		}

		if spec.MatchValue == "" {
			return fmt.Errorf("cant stamp labelspec with empty value")
		}
		m.Labels[spec.MatchLabel] = spec.MatchValue
		return nil
	})
}
