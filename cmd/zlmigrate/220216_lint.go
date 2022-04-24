package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/go-clix/cli"
	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

type linter struct {
	st   zettel.Zetteler
	wg   *sync.WaitGroup
	nags chan<- error
}

func makeCmdLint() *cli.Command {
	cmd := &cli.Command{
		Use:   "lint",
		Short: "warn about markdown-style refs",
	}
	cmd.Run = func(cmd *cli.Command, args []string) error {
		zlpath, ok := os.LookupEnv("ZLPATH")
		if !ok {
			return fmt.Errorf("expected ZLPATH environment variable")
		}

		dir := osfs.New(zlpath)
		store, err := storage.NewStore(dir)
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		nags := make(chan error)
		linter := linter{store, &wg, nags}

		for iter := store.Iter(); iter.Next(); {
			wg.Add(1)
			go linter.lint(iter.Zet())
		}
		go func() {
			wg.Wait()
			close(nags)
		}()

		for nag := range nags {
			log.Println(nag)
		}
		return nil
	}
	return cmd
}

func (l linter) lint(z zettel.Z) {
	defer l.wg.Done()

	scn := bufio.NewScanner(strings.NewReader(z.Readme().Text))
	reg := regexp.MustCompile(`\[.+\]\((.+)\)`)
	for scn.Scan() {
		line := scn.Text()
		matches := reg.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			id := strings.Trim(m[1], " /")
			zl, err := l.st.Zettel(id)
			if err != nil {
				continue
			}

			l.nags <- fmt.Errorf("%q references %q in markdown-link style", z, zl)
		}
	}
}
