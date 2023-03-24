package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"git.jensch.dev/joshua/zl/pkg/zconf"
	"git.jensch.dev/joshua/zl/pkg/zettel"
	"git.jensch.dev/joshua/zl/pkg/zettel/graph"
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/graph/community"
)

//go:embed templates
var tmplFs embed.FS

//go:embed assets/*
var assetFs embed.FS

type server struct {
	templates *template.Template
	engine    *gin.Engine
	store     zettel.Storage
}

func NewServer(store zettel.Storage) (*gin.Engine, error) {
	tmpl, err := template.ParseFS(tmplFs, "templates/*")
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(strings.Split(os.Getenv("ZLSRV_TRUSTED_PROXY"), ","))
	r.SetHTMLTemplate(tmpl)

	server{
		engine:    r,
		templates: tmpl,
		store:     store,
	}.Bind()

	return r, nil
}

func (s server) Bind() {
	s.engine.GET("/", s.root)
	assets, _ := assetFs.ReadDir("assets")
	for _, file := range assets {
		s.engine.GET(fmt.Sprintf("/%s", file.Name()), func(ctx *gin.Context) {
			ctx.FileFromFS(fmt.Sprintf("assets/%s", file.Name()), http.FS(assetFs))
		})
	}
	s.engine.GET("/:zets", s.getFeed)

	api := s.engine.Group("api")
	api.Use(cors)
	api.GET("zettel/:zet", s.apiGetZet)
	api.GET("communities", s.apiGetCommunities)
}

func (s server) root(ctx *gin.Context) {
	entry, ok := os.LookupEnv("ZLSRV_ENTRYPOINT")
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusFound, entry)
}

func (s server) getFeed(ctx *gin.Context) {
	ids := strings.Split(ctx.Param("zets"), ",")
	zets := make(map[string]zettel.Z, len(ids))
	for _, id := range ids {
		zet, err := s.store.Zettel(id)
		if err != nil {
			continue
		}
		zets[zet.Id()] = zet
	}

	conf, err := zconf.FromEnv()
	if err != nil {
		ctx.Error(err)
	}

	g, err := graph.Make(s.store)
	if err != nil {
		ctx.AbortWithError(500, fmt.Errorf("failed creating graph"))
		return
	}

	renderers := make([]ZetRenderer, 0, len(zets))
	base := new(url.URL)
	*base = *ctx.Request.URL
	base.Path = path.Dir(base.Path)
	for _, id := range ids {
		zet, ok := zets[id]
		if !ok {
			continue
		}
		renderers = append(renderers, ZetRenderer{
			Z:       zet,
			G:       g,
			Cfg:     conf,
			Feed:    ids,
			MakeUrl: UrlMaker{base}.MakeUrl,
			Store:   s.store,
			Tmpl:    s.templates,
		})
	}

	ctx.HTML(http.StatusOK, "index.tmpl", renderers)
}

func (s server) apiGetZet(ctx *gin.Context) {
	id := ctx.Param("zet")
	zet, err := s.store.Zettel(id)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}

	data := map[string]interface{}{
		"id": zet.Id(),
		"readme": map[string]string{
			"title": zet.Readme().Title,
			"text":  zet.Readme().Text,
		},
	}
	meta := zet.Metadata()
	if meta != nil {
		mdata := map[string]interface{}{
			"creationTimestamp": meta.CreateTime,
			"labels":            meta.Labels,
		}
		data["meta"] = mdata
	}
	ctx.JSON(http.StatusOK, data)
}

func (s server) apiGetCommunities(ctx *gin.Context) {
	g, err := graph.Make(s.store)
	if err != nil {
		ctx.AbortWithError(500, err)
	}
	reduced := community.Modularize(g, 2, nil)

	comm := reduced.Communities()
	zcmap := make(map[int64]int)
	data := make([]map[string]any, 0, len(comm))
	for i, com := range comm {
		if len(com) < 2 {
			continue
		}
		for _, n := range com {
			zcmap[n.ID()] = i
		}
	}
	for i, com := range comm {
		if len(com) < 2 {
			continue
		}
		comzets := make(map[string]string, len(com))
		comreach := make(map[int]int)
		comdata := map[string]any{
			"id":     i,
			"zettel": comzets,
			"reach":  comreach,
		}
		for _, n := range com {
			z := g.Verts[n.ID()].Z
			comzets[z.Id()] = z.Readme().Title
			edges := g.From(n.ID())
			for edges.Next() {
				// lookup zcmap[nodeid], if it exists, bump that thing up
				comId, ok := zcmap[edges.Node().ID()]
				if !ok || comId == i {
					continue
				}
				comreach[comId] = comreach[comId] + 1
			}
		}
		data = append(data, comdata)
	}
	ctx.PureJSON(http.StatusOK, data)
}
