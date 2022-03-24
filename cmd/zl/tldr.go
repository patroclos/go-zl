package main

import (
	"fmt"
	"strings"

	"github.com/go-clix/cli"
	"github.com/jdkato/prose/summarize"
	"jensch.works/zl/pkg/zettel"
)

func makeCmdTldr(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{
		Use:  "tldr",
		Args: cli.ArgsMin(1),
	}

	x := cmd.Flags().IntP("count", "c", 1, "Number of paragraphs")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		q := strings.Join(args, " ")
		zets, err := st.Resolve(q)
		if err != nil {
			return err
		}

		zl, err := pickOne(zets)
		if err != nil {
			return err
		}

		doc := summarize.NewDocument(zl.Readme().Text)
		tldr := doc.Summary(*x)

		for _, p := range tldr {
			for _, s := range p.Sentences {
				fmt.Println(s.Text)
			}
		}

		return nil
	}
	return cmd
}
