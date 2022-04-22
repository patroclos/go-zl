package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
