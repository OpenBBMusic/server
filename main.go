package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/bb-music/server/internal/api"
	"github.com/bb-music/server/middlewares"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

//go:embed dist/*
var assetFS embed.FS

func main() {
	port := *flag.Int("p", 9799, "启动端口号")
	dev := *flag.Bool("dev", false, "以开发模式启动")
	flag.Parse()

	configDir := GetConfigDir()

	log.Println("configDir", configDir)
	if !dev {
		fmt.Println("生产环境")
		gin.SetMode(gin.ReleaseMode)
		log.SetOutput(&lumberjack.Logger{
			Filename:   filepath.Join(configDir, "logs/log.log"),
			MaxSize:    100,   // 在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxBackups: 10,    // 保留旧文件的最大个数
			MaxAge:     30,    // 保留旧文件的最大天数
			Compress:   false, // 是否压缩/归档旧文件
		})
	} else {
		log.Println("开发环境")
	}

	r := gin.New()
	r.Use(middlewares.Cors(), middlewares.FeAssets(assetFS), middlewares.RequestLogger(), gin.Recovery())

	srv := api.NewServer(r, port, configDir)
	log.Printf("服务已启动：http://127.0.0.1:%v\n", port)
	srv.Run()
}

func GetConfigDir() string {
	configDir, _ := filepath.Abs("./.bb_music")
	return configDir
}
