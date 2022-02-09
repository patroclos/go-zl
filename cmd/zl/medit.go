package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-clix/cli"
	"gopkg.in/yaml.v2"
	"jensch.works/zl/pkg/zettel"
)

func makeCmdMetaEdit(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "metaedit"
	cmd.Aliases = []string{"medit"}
	cmd.Run = func(cmd *cli.Command, args []string) error {
		zets, err := st.Resolve(strings.Join(args, " "))
		if err != nil {
			return err
		}

		zl, err := pickOne(zets)
		if err != nil {
			return err
		}

		tmp, err := zl.Metadata().NewTemp()
		if err != nil {
			return err
		}

		defer os.Remove(tmp.Name())

		shell := exec.Command("vim", tmp.Name())
		shell.Stdin, shell.Stdout, shell.Stderr = os.Stdin, os.Stdout, os.Stderr
		err = shell.Run()
		if _, ok := err.(*exec.ExitError); ok {
			log.Println("aborted")
			return nil
		}

		if err != nil {
			return err
		}

		tmp.Seek(0, io.SeekStart)

		zl, err = zl.Rebuild(func(b zettel.Builder) error {
			dec := yaml.NewDecoder(tmp)
			dec.SetStrict(true)
			m := new(zettel.MetaInfo)
			if err := dec.Decode(m); err != nil {
				return err
			}

			m.CreateTime = b.Metadata().CreateTime
			if b.Metadata().Equal(m) {
				return fmt.Errorf("no change")
			}
			*b.Metadata() = *m

			return nil
		})

		if err != nil {
			return err
		}

		return st.Put(zl)
	}
	return cmd
}
