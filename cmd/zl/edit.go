package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/posener/complete"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/pkg/zettel"
)

type cmdEdit struct {
	ctx *context.Context
}

func (c cmdEdit) Help() string {
	return fmt.Sprintf("Opens a zettel for editing, creating a new git commit")
}
func (c cmdEdit) Run(args []string) int {
	zl, err := c.ctx.Store.Resolve(strings.Join(args, " "))
	if err != nil {
		log.Fatal(err)
	}

	txt, err := ioutil.ReadAll(zl.Reader())
	if err != nil {
		log.Fatal(err)
	}

	tmp, err := os.CreateTemp("", "zledit*.md")
	defer os.Remove(tmp.Name())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(tmp, "# %s\n\n%s", zl.Title(), txt)

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

	if err != nil {
		log.Fatal(err)
	}

	err = c.ctx.Store.Put(zl2)
	if err != nil {
		log.Fatal(err)
	}

	return 0
}
func (c cmdEdit) Synopsis() string {
	return "edit [knode]"
}

func (c cmdEdit) AutocompleteArgs() complete.Predictor {
	iter := c.ctx.Store.Iter()
	set := make([]string, 0, 2048)
	for iter.Next() {
		z := iter.Zet()
		set = append(set, z.Id())
		set = append(set, z.Title())
	}
	return complete.PredictSet(set...)
}

func (c cmdEdit) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}
