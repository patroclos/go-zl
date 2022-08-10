package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-clix/cli"
	"github.com/go-git/go-billy/v5/osfs"
	"git.jensch.dev/zl/pkg/visibility"
	"git.jensch.dev/zl/pkg/zettel"
	"git.jensch.dev/zl/pkg/zettel/elemz"
)

func makeCmdExport(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "export"
	cmd.Short = "Export Listing to OUT directory"
	cmd.Long = "`zl <<< modul | zl export out > out/index.md`"
	cmd.Args = cli.ArgsExact(1)
	tol := cmd.Flags().StringSliceP("tolerate", "t", nil, "comma separated list of tolerated taints")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		if err := os.MkdirAll(args[0], 0700); err != nil {
			return err
		}

		target := osfs.New(args[0])

		scn := elemz.ListScanner(st)

		for zet := range scn.Scan(os.Stdin) {
			if err := target.MkdirAll(zet.Id(), 0700); err != nil {
				return err
			}
			chr, err := target.Chroot(zet.Id())
			if err != nil {
				return err
			}

			if !visibility.Visible(zet, *tol) {
				log.Printf("omitting %s", zet)
				continue
			}

			masked, err := visibility.MaskView{
				Store:    st,
				Tolerate: *tol,
			}.Mask(zet)

			if err != nil {
				return err
			}

			if err := zettel.Write(masked, chr); err != nil {
				return err
			}

			fmt.Println(zettel.MustFmt(masked, zettel.ListingFormat))
		}

		return nil
	}

	return cmd
}
