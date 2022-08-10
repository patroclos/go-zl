package main

import (
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"git.jensch.dev/zl/pkg/zconf"
	"git.jensch.dev/zl/pkg/zettel"
	"git.jensch.dev/zl/pkg/zettel/elemz"
	"git.jensch.dev/zl/pkg/zettel/graph"
)

type ZetRenderer struct {
	Z       zettel.Z
	G       *graph.Graph
	Cfg     *zconf.Cfg
	Feed    []string
	MakeUrl func([]string, *string) *url.URL
	Store   zettel.ZettelerIter
	Tmpl    *template.Template
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
	Rel    string
	InFeed bool
	Type   RefType
}

type RefType string

const (
	RefZ   = "zettel"
	RefUri = "uri"
)

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

	elems, err := elemz.Read(c.Z.Readme().Text)
	if err != nil {
		log.Println("failed reading zet elems", c.Z.Id(), err)
		return
	}
	blinks := c.newBlinksBox()
	if len(blinks.Refs) > 0 {
		elems = append(elems, blinks)
	}

	for _, el := range elems {
		switch elem := el.(type) {
		case *elemz.Refbox:
			c.refbox(*elem)
		case *elemz.Code:
			if c.Cfg.Elems != nil && c.Cfg.Elems.Code != nil {
				flt, ok := c.Cfg.Elems.Code.Filters[elem.BlockParam]
				if !ok {
					goto plain
				}

				cmd := exec.Command("/bin/bash", flt.Cmd)
				cmd.Stdin = strings.NewReader(elem.Code)
				cmd.Stderr = os.Stderr
				buf, err := cmd.Output()
				if err != nil {
					log.Printf("filter error: %v", err)
					goto plain
				}
				_, err = c.sb.Write(buf)
				if err != nil {
					log.Printf("error writing filter-output: %v\nOutput: %#v", err, string(buf))
				}
				continue
			}
		plain:
			c.sb.WriteString(fmt.Sprintf(`<pre class="code"><code>%s</code></pre>`, template.HTMLEscapeString(elem.Code)))
		default:
			c.pre(elem.String())
		}
	}

	return
}

func (c *ZetRenderer) pre(txt string) {
	txt = strings.TrimLeft(txt, "\r\n")
	c.sb.WriteString(fmt.Sprintf("<pre>%s</pre>\n", template.HTMLEscapeString(txt)))
}

func (c ZetRenderer) newBlinksBox() *elemz.Refbox {
	refs := []string{}
	in := c.G.To(c.G.Id(c.Z))
	for in.Next() {
		refs = append(refs, zettel.MustFmt(in.Node().(graph.Node).Z, zettel.ListingFormat))
	}

	l := len(c.Z.Readme().Text)
	return &elemz.Refbox{
		Rel:     "Backlinks",
		Refs:    refs,
		BoxSpan: elemz.Span{Start: l, End: l},
	}
}

func (c *ZetRenderer) refbox(rb elemz.Refbox) {
	data := BoxData{
		Rel:  rb.Rel,
		Refs: make([]BoxRefData, 0, len(rb.Refs)),
	}
	for _, rel := range rb.Refs {
		if strings.HasPrefix(rel, "<") && strings.HasSuffix(rel, ">") {
			rel = rel[1 : len(rel)-1]
			data.Refs = append(data.Refs, BoxRefData{
				Url:    rel,
				Text:   rel,
				InFeed: false,
				Type:   RefUri,
			})
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
		data.Refs = append(data.Refs, BoxRefData{
			Url:    url,
			Text:   refZet.Readme().Title,
			InFeed: hasZet,
			Type:   RefZ,
		})
	}
	if len(data.Refs) == 0 {
		return
	}
	c.Tmpl.ExecuteTemplate(&c.sb, "refbox.tmpl", data)
}

func (c *ZetRenderer) urlTo(id string) string {
	pos := -1
	for i, e := range c.Feed {
		if e == c.Z.Id() {
			pos = i
			break
		}
	}
	if pos == -1 {
		panic("render error: rendered zet not in feed")
	}

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
	for i := range newFeed {
		if i <= pos {
			newFeed[i] = c.Feed[i]
		} else if i == pos+1 {
			newFeed[i] = id
		} else {
			newFeed[i] = c.Feed[i-1]
		}
	}

	return c.MakeUrl(newFeed, &id).String()
}
