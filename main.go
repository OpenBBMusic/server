package main

import (
	"embed"

	"github.com/bb-music/server/internal/server"
	"github.com/bb-music/server/middlewares"
	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var assetFS embed.FS

func main() {
	server.Run(func(r *gin.Engine) {
		r.Use(middlewares.FeAssets(assetFS))
	})
}
