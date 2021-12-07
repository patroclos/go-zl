package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"gopkg.in/yaml.v2"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/storage/filesystem"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

var st storage.Storer

func main() {
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		log.Fatal("no ZLPATH environment present")
	}
	st = &filesystem.ZettelStorage{
		Directory: zlpath,
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/zettel/", zettelHandler)
	http.HandleFunc("/graph", graphHandler)
	// FIXME hardcoded address and tls params
	log.Fatal(http.ListenAndServeTLS(":8087", "crt.pem", "key.pem", nil))
}

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "<html><body><ol>")
	for zl := range storage.AllChan(st) {
		if meta, err := zl.Metadata(); err == nil {
			if _, ok := meta.Labels["zl/inbox"]; !ok {
				continue
			}
		}
		fmt.Fprintf(rw, "<li>")
		rw.Write([]byte(zettel.MustFmt(zl, `<a href="/zettel/{{.Id}}">{{.Title}}</a>`)))
		fmt.Fprintf(rw, "</li>")
	}
	fmt.Fprintf(rw, "</ol></body></html>")
}

func zettelHandler(rw http.ResponseWriter, req *http.Request) {
	pathComps := strings.Split(req.URL.Path, "/")[2:]
	id := pathComps[0]
	zl, err := st.Zettel(zettel.Id(id))
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if len(pathComps) == 2 {
		switch pathComps[1] {
		case "text":
			io.Copy(rw, zl)
			return
		default:
			log.Println("unhandled case, pathcomp not text", pathComps)
		}
	}

	data, err := makeZ(zl)
	if err != nil {
		log.Println(err)
		return
	}
	jsonResp, err := yaml.Marshal(data)
	if err != nil {
		log.Fatalf("Error happened marshaling JSON: %s", err)
	}
	rw.Write(jsonResp)
}
func graphHandler(rw http.ResponseWriter, req *http.Request) {
	g := new(Gra)
	g.Links = make([]*GraLnk, 0, 256)
	g.Nodes = make([]*GraNode, 0, 512)

	nodes := make(map[string]*GraNode)
	for zl := range storage.AllChan(st) {
		id := string(zl.Id())
		origin, ok := nodes[id]
		if !ok {
			origin = NewNode(zl)
			nodes[id] = origin
		}

		if meta, err := zl.Metadata(); err == nil && meta.Link != nil {
			a, b := meta.Link.A, meta.Link.B
			add := func(t string) {
				tn, ok := nodes[t]
				zlr, err := st.Zettel(zettel.Id(t))
				if err != nil {
					return
				}
				if !ok {
					tn = NewNode(zlr)
					nodes[t] = tn
				}
				if ok {
					g.Links = append(g.Links, &GraLnk{Source: origin, Target: tn})
				}
			}

			add(a)
			add(b)
		}
		for zlr := range scan.ListScanner(st).Scan(zl) {
			id2 := string(zlr.Id())
			if id == id2 {
				continue
			}
			target, ok := nodes[id2]
			if !ok {
				target = NewNode(zlr)
				nodes[id2] = target
			}

			lnk := &GraLnk{
				Source: origin,
				Target: target,
			}
			g.Links = append(g.Links, lnk)
		}
	}

	for _, n := range nodes {
		g.Nodes = append(g.Nodes, n)
	}

	log.Printf("Generated graph with %d nodes and %d links", len(g.Nodes), len(g.Links))

	jsonResp, err := json.Marshal(g)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	header := rw.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	rw.Write(jsonResp)
}

type z struct {
	Id      string           `yaml:"id"`
	Title   string           `yaml:"title"`
	Meta    *zettel.MetaInfo `yaml:"meta"`
	TxtHref string           `yaml:"textHref"`
}
type Gra struct {
	Nodes []*GraNode `json:"nodes"`
	Links []*GraLnk  `json:"links"`
}

type GraNode struct {
	Id     string       `json:"id"`
	Title  string       `json:"title"`
	Length int          `json:"length"`
	Meta   *GraNodeMeta `json:"meta"`
}

type GraNodeMeta struct {
	Labels     map[string]string `json:"labels"`
	CreateTime string            `json:"creationTimestamp"`
}

type GraLnk struct {
	Source *GraNode `json:"source"`
	Target *GraNode `json:"target"`
}

func NewNode(zl zettel.Zettel) *GraNode {
	txt, err := zl.Text()
	l := 0
	if err == nil {
		l = utf8.RuneCountInString(txt)

	}
	n := &GraNode{
		Id:     string(zl.Id()),
		Title:  zl.Title(),
		Length: l,
	}
	if meta, err := zl.Metadata(); err == nil {
		n.Meta = &GraNodeMeta{
			Labels:     meta.Labels,
			CreateTime: meta.CreateTime.Format(time.RFC3339),
		}
	}
	return n
}
