package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/search"
	"jensch.works/zl/pkg/zettel"
)

func makeCmdPlace(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "place <z-spec>"
	id := cmd.Flags().String("id", "", `default is newly generated for today`)
	wantEdit := cmd.Flags().BoolP("edit", "e", false, `open the new-Z in $EDITOR`)
	inbox := cmd.Flags().StringP("inbox", "i", "", "set label zl/inbox to the provided value")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		base, err := zettel.Build(func(b zettel.Builder) error {
			switch *id {
			case "":
			default:
				b.Id(*id)
			}

			if *inbox != "" {
				b.Metadata().Labels["zl/inbox"] = *inbox
			}
			q, err := search.Query(strings.Join(args, " "))
			if err != nil {
				return fmt.Errorf("invalid Z-spec: %w", err)
			}

			b.Title(q.Plain)
			for _, spec := range q.Labels {
				if spec.Negated {
					return fmt.Errorf("label-spec can only be positive")
				}
				b.Metadata().Labels[spec.MatchLabel] = spec.MatchValue
			}

			if isTerminal(os.Stdin) {
				return nil
			}

			text, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed reading from stdin: %w", err)
			}

			b.Text(string(text))

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed building Z: %w", err)
		}

		final := base

		if *wantEdit {
			tmp, err := base.Readme().NewTemp()
			if err != nil {
				return fmt.Errorf("error writing Z into temp-file: %w", err)
			}
			defer func() {
				if err != nil {
					fmt.Fprintf(os.Stderr, "unsaved changes: %s\n", tmp.Name())
					return
				}
				os.Remove(tmp.Name())
			}()

			editor, ok := os.LookupEnv("EDITOR")
			if !ok {
				editor = "vim"
			}
			shell := exec.Command(editor, tmp.Name())
			shell.Stdin = os.Stdin
			shell.Stdout = os.Stdout
			shell.Stderr = os.Stderr
			err = shell.Run()
			if _, ok := err.(*exec.ExitError); ok {
				return fmt.Errorf("non-zero exit, aborting")
			}

			if err != nil {
				return fmt.Errorf("editing failed: %w", err)
			}

			tmp.Seek(0, io.SeekStart)
			readme, err := zettel.ParseReadme(tmp)
			if err != nil {
				return fmt.Errorf("failed re-reading edited Z: %w", err)
			}

			final, err = base.Rebuild(func(b zettel.Builder) error {
				b.Title(readme.Title)
				b.Text(readme.Text)
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed assembling Z: %w", err)
			}

		}

		err = st.Put(final)
		if err != nil {
			return fmt.Errorf("failed saving zettel: %w", err)
		}
		fmt.Println(printZet(final))
		return nil
	}

	return cmd
}
