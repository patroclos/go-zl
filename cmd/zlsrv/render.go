package main

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

type ZetRenderer struct {
	Z       zettel.Z
	Feed    []string
	MakeUrl func([]string, *string) *url.URL
	Store   zettel.ZettelerIter
	Tmpl    *template.Template
	blinks  *Backlinks
	sb      strings.Builder
}

type BoxData struct {
	Rel  string
	Refs []BoxRefData
	// Extra []string
}

type BoxRefData struct {
	Url    string
	Text   string
	InFeed bool
}

func (c ZetRenderer) pos() int {
	for i, id := range c.Feed {
		if id == c.Z.Id() {
			return i
		}
	}
	panic("zet not in renderer feed")
}

func (c ZetRenderer) CloseHref() string {
	newFeed := make([]string, 0, len(c.Feed)-1)
	pos := c.pos()

	for i := range c.Feed {
		if i == pos {
			continue
		}
		newFeed = append(newFeed, c.Feed[i])
	}

	if pos <= len(newFeed)-1 {
		return c.MakeUrl(newFeed, &newFeed[pos]).String()
	}

	return c.MakeUrl(newFeed, nil).String()
}

func (c *ZetRenderer) Rendered() (html template.HTML) {
	defer func() {
		html = template.HTML(c.sb.String())
	}()
	if c.sb.Len() > 0 {
		return
	}

	c.sb.WriteString(fmt.Sprintf("<h2>%s</h2>\n", template.HTMLEscapeString(c.Z.Readme().Title)))

	text := c.Z.Readme().Text
	boxes := scan.All(text)
	blinks := c.backlinks()
	if len(blinks.Refs) > 0 {
		boxes = append(boxes, c.backlinks())
	}
	pos := 0

	for _, box := range boxes {
		if box.Start > pos {
			c.pre(text[pos:box.Start])
		}
		c.refbox(box)
		pos = box.End
	}

	if pos < len(text)-1 {
		c.pre(text[pos:])
	}

	return
}

func (c *ZetRenderer) pre(txt string) {
	txt = strings.TrimLeft(txt, "\r\n")
	c.sb.WriteString(fmt.Sprintf("<pre>%s</pre>\n", template.HTMLEscapeString(txt)))
}

func (c ZetRenderer) backlinks() scan.Refbox {
	refs := c.blinks.To(c.Z.Id())

	l := len(c.Z.Readme().Text)
	return scan.Refbox{
		Rel:   "Backlinks",
		Refs:  refs,
		Start: l,
		End:   l,
	}
}

func (c *ZetRenderer) refbox(rb scan.Refbox) {
	data := BoxData{
		Rel:  rb.Rel,
		Refs: make([]BoxRefData, 0, len(rb.Refs)),
	}
	for _, rel := range rb.Refs {
		if strings.HasPrefix(rel, "<") && strings.HasSuffix(rel, ">") {
			rel = rel[1 : len(rel)-1]
			data.Refs = append(data.Refs, BoxRefData{rel, rel, false})
			continue
		}

		refZet, err := c.Store.Zettel(strings.Fields(rel)[0])
		if err != nil {
			continue
		}

		url := c.urlTo(refZet.Id())
		hasZet := false
		for i := range c.Feed {
			if c.Feed[i] == refZet.Id() {
				hasZet = true
				break
			}
		}
		data.Refs = append(data.Refs, BoxRefData{url, refZet.Readme().Title, hasZet})
	}
	if len(data.Refs) == 0 {
		return
	}
	c.Tmpl.ExecuteTemplate(&c.sb, "refbox.tmpl", data)
}

func (c *ZetRenderer) urlTo(id string) string {
	has := false
	for _, entry := range c.Feed {
		if entry == id {
			has = true
			break
		}
	}

	if has {
		return c.MakeUrl(c.Feed, &id).String()
	}

	newFeed := make([]string, len(c.Feed)+1)
	for i := range c.Feed {
		newFeed[i] = c.Feed[i]
	}
	newFeed[len(c.Feed)] = id

	return c.MakeUrl(newFeed, &id).String()
}
