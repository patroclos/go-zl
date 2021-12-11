package zettel

import "fmt"

// error passieren beim machen, nicht beim sein
type Err interface {
	error
	Zettel() Zettel
}

func Errorf(zl Zettel, t string, args ...interface{}) Err {
	return zettelErr{zl, t, args}
}

type zettelErr struct {
	zl   Zettel
	t    string
	args []interface{}
}

func (e zettelErr) Zettel() Zettel { return e.zl }

func (e zettelErr) Error() string {
	txt := fmt.Sprintf(e.t, e.args...)
	idshort := e.zl.Id()
	if len(idshort) > 5 {
		idshort = idshort[len(idshort)-5:]
	}
	return fmt.Sprintf("Err processing Zettel (%s: %s): %s", idshort, e.zl.Title(), txt)
}
