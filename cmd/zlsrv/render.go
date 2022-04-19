package main

import (
	"fmt"
	"html/template"
	"log"
	"net/url"
	"strings"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

type ZetRenderer struct {
	Z       zettel.Zettel
	Feed    []string
	MakeUrl func([]string, *string) *url.URL
	Store   zettel.Zetteler
	Tmpl    *template.Template
	sb      strings.Builder
}

type BoxData struct {
	Rel  string
	Refs []BoxRefData
	// Extra []string
}

type BoxRefData struct {
	Url  string
	Text string
}

func (c *ZetRenderer) Rendered() (html template.HTML) {
	defer func() {
		html = template.HTML(c.sb.String())
	}()
	if c.sb.Len() > 0 {
		return
	}

	c.sb.WriteString(fmt.Sprintf("<h2 id=\"%s\">%s</h2>\n", c.Z.Id(), template.HTMLEscapeString(c.Z.Readme().Title)))

	text := c.Z.Readme().Text
	boxes := scan.All(text)
	pos := 0

	for _, box := range boxes {
		if box.Start > pos {
			c.sb.WriteString(fmt.Sprintf("<pre>%s</pre>\n", template.HTMLEscapeString(text[pos:box.Start])))
		}
		c.refbox(box)
		pos = box.End
	}

	if pos < len(text)-1 {
		c.sb.WriteString(fmt.Sprintf("<pre>%s</pre>\n", template.HTMLEscapeString(text[pos:])))
	}

	return
}

func (c *ZetRenderer) refbox(rb scan.Refbox) {
	data := BoxData{
		Rel:  rb.Rel,
		Refs: make([]BoxRefData, 0, len(rb.Refs)),
	}
	for _, rel := range rb.Refs {
		if strings.HasPrefix(rel, "<") && strings.HasSuffix(rel, ">") {
			rel = rel[1 : len(rel)-1]
			data.Refs = append(data.Refs, BoxRefData{rel, rel})
			continue
		}

		refZet, err := c.Store.Zettel(strings.Fields(rel)[0])
		if err != nil {
			log.Println(fmt.Errorf("failed resolving refbox ref %q: %w", rel, err))
			continue
		}

		url := c.urlTo(refZet.Id())
		data.Refs = append(data.Refs, BoxRefData{url, refZet.Readme().Title})
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
