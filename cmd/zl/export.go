package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/visibility"
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

			zet, err = zet.Rebuild(func(b zettel.Builder) error {
				var str strings.Builder
				elems, err := scan.Elements(st, zet.Readme().Text)
				if err != nil {
					return err
				}
				lines := make([]string, 0)
				scn := bufio.NewScanner(strings.NewReader(zet.Readme().Text))
				for scn.Scan() {
					lines = append(lines, scn.Text())
				}
				for i := range elems {
					e := elems[i]
					var txt string
					if e.Span.Start == -1 {
						txt = e.Span.Input
					} else {
						txt = strings.Join(lines[e.Span.Start:e.Span.Pos], "\n")
					}
					if e.Type != scan.ItemRefbox {
						log.Println(e.Type, e.Span)
						str.WriteString(fmt.Sprintln(txt))
						continue
					}
					scn := bufio.NewScanner(strings.NewReader(txt))
					scn.Scan()
					str.WriteString(fmt.Sprintln(scn.Text()))
					for scn.Scan() {
						line := scn.Text()
						zets, err := st.Resolve(line)
						if err != nil {
							return fmt.Errorf("refbox entry not resolved: %w", err)
						}
						if len(zets) != 1 {
							return fmt.Errorf("refbox entry %q must be unique, but matches: %v", line, zets)
						}
						z := zets[0]
						tolerate := strings.Split(os.Getenv("ZL_TOLERATE"), ",")
						log.Println("viz?")
						if !visibility.Visible(z, tolerate) {
							log.Println("not viz")
							z, err = z.Rebuild(func(b zettel.Builder) error {
								b.Title("MASKED")
								b.Text("")
								return nil
							})
							if err != nil {
								return err
							}
						}
						str.WriteString(fmt.Sprintln(z))
					}
				}
				b.Text(str.String())
				return nil
			})

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
