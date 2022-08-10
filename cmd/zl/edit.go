package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-clix/cli"
	"git.jensch.dev/zl/pkg/zettel"
)

func makeCmdEdit(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{}
	cmd.Use = "edit zetref"
	cmd.Run = func(cmd *cli.Command, args []string) error {
		zets, err := st.Resolve(strings.Join(args, " "))
		if err != nil {
			return err
		}
		zl, err := pickOne(zets)
		if err != nil {
			return err
		}

		tmp, err := zl.Readme().NewTemp()
		if err != nil {
			return err
		}

		// keep in mind, the following log.Fatal calls will circumvent this
		defer os.Remove(tmp.Name())

		shell := exec.Command("vim", tmp.Name())
		shell.Stdin = os.Stdin
		shell.Stdout = os.Stdout
		shell.Stderr = os.Stderr
		err = shell.Run()
		if _, ok := err.(*exec.ExitError); ok {
			log.Println("aborted")
			return nil
		}

		if err != nil {
			return err
		}

		tmp.Seek(0, io.SeekStart)
		readme, err := zettel.ParseReadme(tmp)
		if err != nil {
			return err
		}

		if *readme == zl.Readme() {
			log.Println("nothing changed")
			return nil
		}

		zl2, err := zl.Rebuild(func(b zettel.Builder) error {
			b.Title(readme.Title)
			b.Text(readme.Text)
			return nil
		})

		if err != nil {
			return err
		}

		err = st.Put(zl2)
		if err != nil {
			return err
		}

		return nil
	}
	return cmd
}
