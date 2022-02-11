package main

import (
	"fmt"
	"os"

	"github.com/go-clix/cli"
	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

func makeCmdExport(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "export"
	cmd.Short = "Export Listing to OUT directory"
	cmd.Long = "`zl <<< modul | zl export out > out/index.md`"
	cmd.Args = cli.ArgsExact(1)
	cmd.Run = func(cmd *cli.Command, args []string) error {
		if err := os.MkdirAll(args[0], 0700); err != nil {
			return err
		}

		target := osfs.New(args[0])

		scn := scan.ListScanner(st)

		for zet := range scn.Scan(os.Stdin) {
			if err := target.MkdirAll(zet.Id(), 0700); err != nil {
				return err
			}
			chr, err := target.Chroot(zet.Id())
			if err != nil {
				return err
			}

			if err := zettel.Write(zet, chr); err != nil {
				return err
			}

			fmt.Println(zettel.MustFmt(zet, zettel.ListingFormat))
		}

		return nil
	}

	return cmd
}

/*
type cmdExport struct {
	st     zettel.Storage
	target billy.Filesystem
}

func (ex cmdExport) Help() string {
	return ``
}

func (ex cmdExport) Synopsis() string {
	return ``
}

func (ex cmdExport) Run(args []string) int {
	list := args[0]

	if err := os.MkdirAll(args[1], 0700); err != nil {
		return 1
	}

	ex.target = osfs.New(args[1])

	scn := scan.ListScanner(ex.st)
	f, err := os.Open(list)
	if err != nil {
		return 1
	}

	for zet := range scn.Scan(f) {
		if err := ex.target.MkdirAll(zet.Id(), 0700); err != nil {
			return 1
		}
		chr, err := ex.target.Chroot(zet.Id())
		if err != nil {
			return 1
		}

		zettel.Write(zet, chr)
	}

	return 0
}
*/
