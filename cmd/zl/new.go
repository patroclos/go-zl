package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
)

func makeCmdNew(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "new"
	cmd.Run = func(_ *cli.Command, args []string) error {
		title := strings.Join(args, " ")
		zl, err := zettel.Build(func(b zettel.Builder) error {
			b.Title(title)
			b.Metadata().Labels["zl/inbox"] = "default"
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
			return err
		}

		tmp.Seek(0, io.SeekStart)
		readme, err := zettel.ParseReadme(tmp)
		if err != nil {
			return err
		}

		zl2, err := zl.Rebuild(func(b zettel.Builder) error {
			b.Title(readme.Title)
			b.Text(readme.Text)
			return nil
		})

		err = st.Put(zl2)
		if err != nil {
			return err
		}

		return nil
	}
	return cmd
}
