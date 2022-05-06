package elemz

type Text struct {
	Content string
	span    Span
}

func (e *Text) Span() Span     { return e.span }
func (e *Text) String() string { return e.Content }
