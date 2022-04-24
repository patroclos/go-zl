package main

import (
	"embed"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"jensch.works/zl/pkg/zettel"
)

//go:embed templates
var tmplFs embed.FS

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
	r.SetTrustedProxies(nil)
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
	s.engine.GET("/:zets", s.getFeed)

	api := s.engine.Group("api")
	api.Use(cors)
	api.GET("zettel/:zet", s.apiGetZet)
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
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		zets[zet.Id()] = zet
	}

	renderers := make([]ZetRenderer, 0, len(zets))
	base := new(url.URL)
	*base = *ctx.Request.URL
	base.Path = path.Dir(base.Path)
	for _, id := range ids {
		renderers = append(renderers, ZetRenderer{
			Z:       zets[id],
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
