package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
)

func main() {
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		log.Fatal("ZLPATH evnironment variable is mandatory.")
	}

	dir := osfs.New(zlpath)
	store, err := storage.NewStore(dir)
	if err != nil {
		log.Fatal(err)
	}

	srv := server{st: store}
	srv.Bind()
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Println(err)
	}
}

type server struct {
	st zettel.Storage
}

func (s server) Bind() {
	http.HandleFunc("/", s.root)
	http.HandleFunc("/zet/", s.zet)
}

func (s server) root(rw http.ResponseWriter, req *http.Request) {
	txt := `
<html>
<head>
<title>zetsrv</title>
</head>
<body>
<style>
* {
background-color: black;
color: white;
}
</style>
<ul>
{{range $z := .}}
<li><a href="/zet/{{$z.Id}}">{{$z.Id}}  {{$z.Readme.Title}}</a></li>
{{end}}
</ul>
</body>
</html>
	`
	tmpl, err := template.New("root").Parse(txt)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	zets := make([]zettel.Zettel, 0, 1024)
	for it := s.st.Iter(); it.Next(); {
		zets = append(zets, it.Zet())
	}
	tmpl.Execute(rw, zets)
}

type zet struct {
	Z   zettel.Zettel
	In  []zettel.Zettel
	Out []zettel.Zettel
}

func (s server) zet(rw http.ResponseWriter, req *http.Request) {
	comps := strings.Split(req.URL.Path, "/")[2:]
	id := comps[0]
	zl, err := s.st.Zettel(id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	in, out := make([]zettel.Zettel, 0, 32), make([]zettel.Zettel, 0, 8)

	crawl.New(s.st, func(n crawl.Node) crawl.RecurseMask {
		if len(n.Path) == 0 {
			return crawl.All
		}

		if n.Reason&crawl.Inbound != 0 {
			in = append(in, n.Z)
		}
		if n.Reason&crawl.Outbound != 0 {
			out = append(out, n.Z)
		}
		return crawl.None
	}).Crawl(zl)

	txt := `
<html>
<head>
<title>{{.Z.Readme.Title}}</title>
</head>
<body>
<style>
* {
background-color: black;
color: white;
text-decoration: none;
}
</style>
<pre>{{.Z.Readme}}</pre>
<h2>Inbound</h2>
<ul>
{{ range $z := .In }}
<li><a href="/zet/{{$z.Id}}">{{$z.Id}}  {{$z.Readme.Title}}</a></li>
{{ end }}
</ul>
<h2>Outbound</h2>
{{ range $z := .Out }}
<li><a href="/zet/{{$z.Id}}">{{$z.Id}}  {{$z.Readme.Title}}</a></li>
{{ end }}
</body>
</html>
	`
	tmpl, err := template.New("zet").Parse(txt)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	log.Println(in, out)
	if err := tmpl.Execute(rw, zet{Z: zl, In: in, Out: out}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
}
