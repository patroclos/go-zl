package view

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"jensch.works/zl/pkg/zettel"
)

func viewZettel(zl zettel.Zettel) error {
	temp, err := os.CreateTemp("", "zl*.md")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	txt, err := zl.Text()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(temp, "# %s\n\n%s", zl.Title(), txt)
	cmd := exec.Command("vim", temp.Name())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err = cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	}

	fmt.Println(zettel.MustFmt(zl, zettel.ListFormat))

	return nil
}

func listFmt(zl zettel.Zettel) string {
	txt, err := zettel.Fmt(zl, zettel.ListFormat)
	if err != nil {
		panic(err)
	}

	return txt
}

func pickOne(from []zettel.Zettel) (zettel.Zettel, error) {
	switch len(from) {
	case 0:
		return nil, fmt.Errorf("tried picking from empty zettel slice")
	case 1:
		return from[0], nil
	default:
		for i, zl := range from {
			fmt.Fprintf(os.Stderr, "[%03d] %s\n", i+1, zettel.MustFmt(zl, zettel.ListFormat))

		}
		fmt.Fprintf(os.Stderr, "Choice: ")
		choice := 0

		for {
			_, err := fmt.Scanf("%d", &choice)

			if err != nil {
				return nil, err
			}
			if choice <= 0 || choice > len(from) {
				log.Printf("%0d is out of range", choice)
				continue
			}

			return from[choice-1], nil
		}
	}
}
