package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

func main() {
	root := cli.Command{
		Use:  "zlmigrate",
		Args: cli.ArgsExact(2),
	}

	root.Run = func(cmd *cli.Command, args []string) error {
		src, err := storage.NewStore(osfs.New(args[0]))
		if err != nil {
			return err
		}
		dst, err := storage.NewStore(osfs.New(args[1]))
		if err != nil {
			return err
		}

		counterparts := make(map[string]zettel.Zettel)
		olds := make([]string, 0, 1024)

		scn := scan.ListScanner(src)

		for iter := src.Iter(); iter.Next(); {
			zet := iter.Zet()
			split := strings.Split(zet.Id(), "-")
			if len(split) != 3 {
				log.Printf("ignoring %s", zet.Id())
				continue
			}
			trimmed, err := zet.Rebuild(func(b zettel.Builder) error {
				newId := fmt.Sprintf("%s-%s", split[0][2:], split[2][:4])
				b.Id(newId)
				return nil
			})
			counterparts[trimmed.Id()] = zet
			counterparts[zet.Id()] = trimmed
			if err != nil {
				return err
			}

			if _, err := dst.Zettel(trimmed.Id()); err == nil {
				continue
			}
			olds = append(olds, zet.Id())

			log.Printf("%s => %s", zet.Id(), trimmed.Id())
		}

		for _, id := range olds {
			new := counterparts[id]
			r := strings.NewReader(new.Readme().Text)
			refs := make([]zettel.Zettel, 0, 8)
			for ref := range scn.Scan(r) {
				refs = append(refs, ref)
			}

			txt := new.Readme().Text
			for _, ref := range refs {
				newRef := counterparts[ref.Id()]
				txt = strings.ReplaceAll(txt, ref.Id(), newRef.Id())
			}

			if len(refs) > 0 {
				fmt.Printf("\n--- %s  %s ---\n%s\n", new.Id(), new.Readme().Title, txt)
			}
			new, err = new.Rebuild(func(b zettel.Builder) error {
				b.Text(txt)
				return nil
			})

			if err != nil {
				return err
			}

			dst.Put(new)
		}

		log.Printf("from: %v\nto: %v\n", src, dst)
		return nil
	}
	if err := root.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
