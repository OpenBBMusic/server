package middlewares

import (
	"embed"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// h5 history api
func FeAssets(assetFS *embed.FS) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if assetFS == nil {
			ctx.Next()
			return
		}
		path := ctx.Request.URL.Path
		// fmt.Println("PATH", path)
		if strings.HasPrefix(path, "/api/") {
			ctx.Next()
			return
		}
		if path == "/" {
			path = "/index.html"
		}
		file := filepath.Join("dist", path)
		data, err := assetFS.ReadFile(file)
		if err != nil {
			ctx.Next()
			return
		}
		httpType := "text/html"
		if filepath.Ext(file) == ".js" {
			httpType = "application/javascript"
		}
		if filepath.Ext(file) == ".css" {
			httpType = "text/css"
		}
		// fmt.Println("httpType", httpType, "data", string(data))
		ctx.Data(http.StatusOK, httpType+"; charset=utf-8", data)
	}
}
