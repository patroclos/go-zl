package lint

import (
	"sync"

	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

type Linter interface {
	Lint()
}

type linter struct {
	It  storage.ZettelIter
	Nag []Nagger
}

func (l linter) Lint() <-chan Nag {
	ch := make(chan Nag)
	go func() {
		wg := new(sync.WaitGroup)
		wg.Add(len(l.Nag))
		for _, n := range l.Nag {
			nagger := n
			go func() {
				defer wg.Done()
				for zl := range storage.AllChan(l.It) {
					for nag := range nagger.Nag(zl) {
						ch <- nag
					}
				}
			}()
		}

		defer close(ch)

		wg.Wait()
	}()
	return ch
}

type Nagger interface {
	Nag(z zettel.Zettel) <-chan Nag
}

type NagLevel string

var (
	Error   = NagLevel("ERR")
	Warn    = NagLevel("WRN")
	Suggest = NagLevel("Suggest")
	Hint    = NagLevel("Hint")
)

type Nag struct {
	Message string
	Level   NagLevel
	Z       zettel.Id
}
