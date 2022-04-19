package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

type UrlMaker struct {
	Base *url.URL
}

func (x UrlMaker) MakeUrl(feed []string, focus *string) *url.URL {
	if focus == nil {
		url, _ := x.Base.Parse(fmt.Sprintf("%s", strings.Join(feed, ",")))
		return url
	}

	url, _ := x.Base.Parse(fmt.Sprintf("%s#%s", strings.Join(feed, ","), *focus))
	return url
}

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

	r := gin.Default()
	r.SetTrustedProxies(nil)

	tmpl, err := template.ParseGlob("templates/*")
	if err != nil {
		log.Fatal(err)
	}
	r.SetHTMLTemplate(tmpl)

	r.GET("/:zets", func(ctx *gin.Context) {
		ids := strings.Split(ctx.Param("zets"), ",")
		zets := make(map[string]zettel.Zettel, len(ids))
		for _, id := range ids {
			zet, err := store.Zettel(id)
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
				Store:   store,
				Tmpl:    tmpl,
			})
		}

		ctx.HTML(http.StatusOK, "index.tmpl", renderers)
	})

	api := r.Group("api")
	api.Use(cors)
	api.GET("zettel/:zet", func(ctx *gin.Context) {
		id := ctx.Param("zet")
		zet, err := store.Zettel(id)
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
	})

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}

func cors(ctx *gin.Context) {
	origin := ctx.Request.Header.Get("Origin")
	if len(origin) == 0 {
		return
	}
	host := ctx.Request.Host

	if origin == "http://"+host || origin == "https://"+host {
		return
	}

	header := ctx.Writer.Header()
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Headers", "*")
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
	} else {
		ctx.Next()
	}
}
