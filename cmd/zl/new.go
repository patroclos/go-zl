package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/pkg/zettel"
)

type cmdNew struct {
	ctx *context.Context
}

func (c cmdNew) Help() string {
	return "blub"
}

func (c cmdNew) Synopsis() string {
	return "new [title]"
}

func (c cmdNew) Run(args []string) int {
	title := strings.Join(args, " ")
	zl, err := zettel.Build(func(b zettel.Builder) error {
		b.Title(title)
		return nil
	})
	if err != nil {
		log.Fatalf("failed creating Zettel: %v", err)
	}

	tmp, err := zl.Readme().NewTemp()
	if err != nil {
		log.Fatalf("failed creating tmp file: %v", err)
	}

	defer tmp.Close()

	cmd := exec.Command("vim", tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		log.Fatalln("aborted")
	}

	if err != nil {
		log.Fatal(err)
	}

	tmp.Seek(0, io.SeekStart)
	readme, err := zettel.ParseReadme(tmp)
	if err != nil {
		log.Fatal(err)
	}

	zl2, err := zl.Rebuild(func(b zettel.Builder) error {
		b.Title(readme.Title)
		b.Text(readme.Text)
		return nil
	})

	err = c.ctx.Store.Put(zl2)
	if err != nil {
		log.Fatal(err)
	}

	return 0
}
