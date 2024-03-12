package main

import (
	"flag"
	"os"

	pullfe "github.com/bb-music/server/internal/pull_fe"
	"github.com/bb-music/server/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {
	token := flag.String("token", "", "Your Github token")
	flag.Parse()

	var t = *token

	if *token == "" {
		t = os.Getenv("GITHUB_TOKEN")
	}

	pullfe.Start(t)
	server.Run(func(r *gin.Engine) {})
}
