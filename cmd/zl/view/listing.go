package view

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"jensch.works/zl/pkg/zettel"
)

type Listing struct {
	Zets  []zettel.Z
	Fmt   string
	Dest  io.Writer
	Color bool
}

func (v *Listing) Render() error {

	for _, z := range v.Zets {
		format := v.Fmt
		switch format {
		case "listing":
			if !isTerminal(v.Dest) {
				format = zettel.ListingFormat
				break
			}
			var fmt strings.Builder
			gray := "\x1b[38;5;242m"
			reset := "\x1b[0m"

			inbox, hasInbox := z.Metadata().Labels["zl/inbox"]
			hasInbox = hasInbox && v.Color
			if hasInbox {
				fmt.WriteString(gray)
			}
			fmt.WriteString("{{.Id}}")
			if hasInbox {
				fmt.WriteString(reset)
			}
			fmt.WriteString("  {{.Title}}")

			if hasInbox && inbox != "default" {
				fmt.WriteString("  ")
				lipgloss.NewStyle().Foreground(lipgloss.Color("gray")).Render(`l:{{index .Labels "zl/inbox}}"`)
			}
			format = fmt.String()
		case "list":
			format = zettel.ListFormat
		}

		txt, err := zettel.Fmt(z, format)
		if err != nil {
			fmt.Printf("%s", z.Id())
			// return fmt.Errorf("formatting zet %s %q failed: %w", z.Id(), format, err)
		}

		fmt.Fprintln(v.Dest, txt)
	}
	return nil
}

func isTerminal(f io.Writer) bool {
	type Stater interface {
		Stat() (os.FileInfo, error)
	}

	if f == nil {
		return false
	}
	stat, ok := f.(Stater)
	if !ok {
		return false
	}

	o, err := stat.Stat()
	if err != nil {
		return false
	}

	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}
