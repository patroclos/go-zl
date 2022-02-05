package main

import (
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

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
