package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
			migrated, err := zet.Rebuild(func(b zettel.Builder) error {
				newId := fmt.Sprintf("%s-%s", split[0][2:], split[2][:4])
				b.Id(newId)

				meta := b.Metadata()
				t, err := time.Parse("20060102150405", fmt.Sprintf("%s%s", split[0], split[1]))
				if err != nil {
					log.Fatal(fmt.Errorf("failed parsing timestamp %s: %w", zet.Id(), err))
				}
				meta.CreateTime = t

				if meta.Link != nil {
					if a, ok := counterparts[meta.Link.A]; ok {
						meta.Link.A = a.Id()
					}
					if b, ok := counterparts[meta.Link.B]; ok {
						meta.Link.B = b.Id()
					}
					for i := range meta.Link.Ctx {
						if c, ok := counterparts[meta.Link.Ctx[i]]; ok {
							meta.Link.Ctx[i] = c.Id()
						}
					}
				}

				return nil
			})
			counterparts[migrated.Id()] = zet
			counterparts[zet.Id()] = migrated
			if err != nil {
				return err
			}

			if _, err := dst.Zettel(migrated.Id()); err == nil {
				continue
			}
			olds = append(olds, zet.Id())

			log.Printf("%s => %s", zet.Id(), migrated.Id())
		}

		for _, id := range olds {
			new := counterparts[id]
			r := strings.NewReader(new.Readme().Text)
			refs := make([]zettel.Zettel, 0, 8)
			for ref := range scn.Scan(r) {
				refs = append(refs, ref)
			}

			txt := new.Readme().Text
			oldnew := make([]string, len(refs)*2)
			for i, ref := range refs {
				newRef := counterparts[ref.Id()]
				oldnew[i*2] = ref.Id()
				oldnew[i*2+1] = newRef.Id()
				// txt = strings.ReplaceAll(txt, ref.Id(), newRef.Id())
			}

			txt = strings.NewReplacer(oldnew...).Replace(txt)

			if l := len(refs); l > 0 {
				fmt.Printf("rewriting %d references in %s (formerly %s)\n", l, new.Id(), id)
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
