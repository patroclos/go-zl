package scan

import (
	"bufio"
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
	itemErr ItemType = iota
	itemTxt
	itemEof
	itemRefbox
)

func Elements(st zettel.Zetteler, txt string) ([]*Item, error) {
	_, ch := lex(st, txt)
	items := make([]*Item, 0)
	for x := range ch {
		items = append(items, x)
	}
	return items, nil
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	scn   *bufio.Scanner
	st    zettel.Zetteler
	start int
	pos   int
	items chan *Item
}

func lex(st zettel.Zetteler, input string) (*lexer, chan *Item) {
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
		BlkSpan{l.input, l.start, l.pos},
		itemErr,
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
				l.emit(itemTxt)
				l.start = l.pos
				l.emit(itemEof)
				return nil
			}
			potential := l.scn.Text()
			potential = strings.TrimLeft(potential, "* ")
			_, err := l.st.Zettel(strings.Fields(potential)[0])
			if err != nil {
				l.pos += 2
				return lexText
			}
			l.emit(itemTxt)
			l.start, l.pos = l.pos, l.pos+2
			return lexRefBlock
		}
		l.pos += 1
	}
	if l.pos > l.start {
		l.emit(itemTxt)
		l.start = l.pos
	}
	l.emit(itemEof)
	return nil
}

func lexRefBlock(l *lexer) stateFn {
	for l.scn.Scan() {
		line := l.scn.Text()
		if len(line) == 0 {
			l.emit(itemRefbox)
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
				l.emit(itemRefbox)
				l.pos += 1
				return lexText
			}
			if l.pos > l.start {
				l.emit(itemRefbox)
				l.start = l.pos
				l.emit(itemEof)
			}
			return nil
		}
		_, err := l.st.Zettel(l.scn.Text())
		if err != nil {
			l.errorf("%v", err)
			return nil
		}
	}
	if l.start < l.pos {
		l.emit(itemRefbox)
		l.start = l.pos
	}
	return lexText
}
