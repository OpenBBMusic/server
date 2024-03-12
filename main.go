package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strconv"

	pullfe "github.com/bb-music/server/internal/pull_fe"
	"github.com/bb-music/server/internal/server"
)

//go:embed dist/*
var assetFS embed.FS

func main() {
	config := GetConfig()
	fmt.Printf("CONFIG: %+v\n", config)

	var asset *embed.FS
	if config.BuildFeAssets {
		asset = &assetFS
		pullfe.Start(config.GithubToken)
	}

	server.Run(config.Port, config.IsDev, asset)
}

type Config struct {
	GithubToken   string
	BuildFeAssets bool
	IsDev         bool
	Port          int
}

func GetConfig() Config {
	token := flag.String("token", "", "Github Token 用于拉取前端资源")
	buildFe := flag.Bool("build-fe", false, "是否将前端资源打包到二进制中")
	dev := flag.Bool("dev", false, "是否以开发模式启动")
	serverPort := flag.Int("p", 9799, "启动端口号")
	useEnv := flag.Bool("use-env", true, "优先使用环境变量")

	flag.Parse()

	config := Config{
		GithubToken:   *token,
		BuildFeAssets: *buildFe,
		IsDev:         *dev,
		Port:          *serverPort,
	}

	if *useEnv {
		gt := os.Getenv("GITHUB_TOKEN")
		build_fe := os.Getenv("BUILD_FE")
		port := os.Getenv("PORT")
		if gt != "" {
			config.GithubToken = gt
		}
		if build_fe == "1" {
			config.BuildFeAssets = true
		}
		if port != "" {
			// port 字符串转 int
			num, err := strconv.Atoi(port)
			if err == nil {
				config.Port = num
			}
		}
	}

	return config
}
