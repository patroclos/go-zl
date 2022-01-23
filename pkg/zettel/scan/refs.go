package scan

import (
	"log"
	"regexp"
	"strings"
	"sync"

	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

func Refs(text string) []string {
	reg := regexp.MustCompile(`\[.+\]\((.+)\)`)
	matches := reg.FindAllStringSubmatch(text, -1)
	results := make([]string, 0, 8)
	for _, m := range matches {
		id := strings.Trim(m[1], " /")
		results = append(results, id)
	}
	reg = regexp.MustCompile(`\* ([a-zA-Z0-9-]+)  .*`)
	matches = reg.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		id := m[1]
		results = append(results, id)
	}

	return results
}

func Backrefs(to string, st interface {
	storage.ZettelIter
	storage.Zetteler
}) <-chan RefZet {
	ch := make(chan RefZet)
	wg := new(sync.WaitGroup)
	for zl := range storage.AllChan(st) {
		zet := zl
		wg.Add(1)
		go func() {
			defer wg.Done()
			for zlr := range ListScanner(st).Scan(zet.Reader()) {
				if string(zlr.Id()) == string(to) {
					ch <- RefZet{
						Zettel: zet,
						txt:    "<pending>",
						target: zlr,
						link:   nil,
					}
				}
			}

			meta := zet.Metadata()
			if meta.Link == nil {
				return
			}
			lnk := meta.Link
			switch string(to) {
			case lnk.A:
				b, errB := st.Zettel(lnk.B)
				if errB != nil {
					log.Printf("failed opening link.to zettel (%s) from %s: %v", lnk.B, zet.Id(), errB)
					return
				}
				ch <- RefZet{
					Zettel: zet,
					txt:    zet.Title(),
					target: b,
				}

			case lnk.B:
				a, errA := st.Zettel(lnk.A)
				if errA != nil {
					log.Printf("failed opening link.from zettel (%s) from %s: %v", lnk.A, zet.Id(), errA)
					return
				}
				ch <- RefZet{
					Zettel: zet,
					txt:    zet.Title(),
					target: a,
				}
			default:
				return
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

type RefZet struct {
	zettel.Zettel
	txt    string
	target zettel.Zettel
	link   zettel.Zettel
}

// link may be nil, target should always be valid
func (rz RefZet) Ref() (txt string, target zettel.Zettel, link zettel.Zettel) {
	return rz.txt, rz.target, rz.link
}
