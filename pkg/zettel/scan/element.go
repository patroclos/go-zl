package scan

import (
	"bufio"
	"fmt"
	"strings"

	"jensch.works/zl/pkg/zettel"
)

type RefBox struct {
	Span  BlkSpan
	Refs  []zettel.Zettel
	Extra BlkSpan
}

type BlkSpan struct {
	Input      string
	Start, Pos int
}

type Item struct {
	Span BlkSpan
	Type ItemType
}

type ItemType int

const (
	ItemErr ItemType = iota
	ItemTxt
	itemEof
	ItemRefbox
)

func Elements(st zettel.Resolver, txt string) ([]*Item, error) {
	_, ch := lex(st, txt)
	items := make([]*Item, 0)
	for x := range ch {
		if x.Type == ItemErr {
			return nil, fmt.Errorf("%sl", x.Span.Input)
		}
		items = append(items, x)
	}
	return items, nil
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	scn   *bufio.Scanner
	st    zettel.Resolver
	start int
	pos   int
	items chan *Item
}

func lex(st zettel.Resolver, input string) (*lexer, chan *Item) {
	l := &lexer{
		st:    st,
		items: make(chan *Item),
		input: input,
		scn:   bufio.NewScanner(strings.NewReader(input)),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- &Item{
		BlkSpan{fmt.Sprintf(format, args...), -1, -1},
		ItemErr,
	}
	return nil
}
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t ItemType) {
	l.items <- &Item{BlkSpan{l.input, l.start, l.pos}, t}
	l.start = l.pos
}

func lexText(l *lexer) stateFn {
	for l.scn.Scan() {
		line := l.scn.Text()
		if strings.HasSuffix(line, ":") {
			if !l.scn.Scan() {
				l.emit(ItemTxt)
				l.start = l.pos
				l.emit(itemEof)
				return nil
			}
			potential := l.scn.Text()
			_, err := l.st.Resolve(potential)
			if err != nil {
				l.pos += 2
				return lexText
			}
			l.emit(ItemTxt)
			l.start, l.pos = l.pos, l.pos+2
			return lexRefBlock
		}
		l.pos += 1
	}
	if l.pos > l.start {
		l.emit(ItemTxt)
		l.start = l.pos
	}
	l.emit(itemEof)
	return nil
}

func lexRefBlock(l *lexer) stateFn {
	for l.scn.Scan() {
		line := l.scn.Text()
		if len(line) == 0 {
			l.emit(ItemRefbox)
			l.start = l.pos
			return lexText
		}
		l.pos += 1
		if strings.HasPrefix(line, "+ ") {
			for l.scn.Scan() {
				ln := l.scn.Text()
				if strings.HasPrefix(ln, "  ") {
					l.pos += 1
					continue
				}
				l.emit(ItemRefbox)
				l.pos += 1
				return lexText
			}
			if l.pos > l.start {
				l.emit(ItemRefbox)
				l.start = l.pos
				l.emit(itemEof)
			}
			return nil
		}
		_, err := l.st.Resolve(l.scn.Text())
		if err != nil {
			l.errorf("%v", err)
			return nil
		}
	}
	if l.start < l.pos {
		l.emit(ItemRefbox)
		l.start = l.pos
	}
	return lexText
}
