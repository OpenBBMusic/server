package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/bb-music/desktop/app_bili"
	"github.com/bb-music/desktop/pkg/bb_server"
	"github.com/bb-music/desktop/pkg/bb_type"
	"github.com/bb-music/server/middlewares"
	"github.com/bb-music/server/resp"
	"github.com/bb-music/server/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	configDir := GetConfigDir()
	log.Println("configDir", configDir)
	if !utils.IsDev() {
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

	port := 9799
	srv := NewServer(port, configDir)
	log.Printf("服务已启动：http://127.0.0.1:%v\n", port)
	srv.Run()
}

func GetConfigDir() string {
	configDir, _ := filepath.Abs("./.bb_music")
	return configDir
}

type OriginService interface {
	GetConfig() any
	InitConfig(force bool) error
	Search(params bb_type.SearchParams) (*bb_type.SearchResponse, error)
	SearchDetail(id string) (*bb_type.SearchItem, error)
	GetMusicFile(id string) (*httputil.ReverseProxy, *http.Request, error)
	DownloadMusic(params bb_type.DownloadMusicParams) (string, error)
}

func NewServer(port int, cacheDir string) *bb_server.Server {
	log.Println("======= 服务启动 =======")
	log.Println("端口:", port)
	log.Println("缓存目录:", cacheDir)
	log.Println("注册音乐源服务")
	// 日志输出
	svcLogger := NewSvcLogger()

	bili := app_bili.New(cacheDir, svcLogger)
	// 注册源服务
	service := map[bb_type.OriginType]OriginService{
		bb_type.BiliOriginName: bili,
	}
	bili.InitConfig(false)

	// 初始化 gin
	r := gin.New()
	r.Use(middlewares.Cors(), middlewares.RequestLogger(), gin.Recovery())

	// 获取源的配置信息
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "欢迎使用 哔哔音乐 API 服务")
	})

	g := r.Group("/api")

	// 获取源的配置信息
	g.GET("/config/:origin", func(ctx *gin.Context) {
		origin := ctx.Param("origin")
		data := service[origin].GetConfig()
		ctx.JSON(resp.Success(data, "查询成功"))
	})

	// 搜索音乐
	g.GET("/search/:origin", func(ctx *gin.Context) {
		origin := ctx.Param("origin")
		keyword := ctx.Query("keyword")
		page := ctx.DefaultQuery("page", "1")

		data, err := service[origin].Search(bb_type.SearchParams{
			Keyword: keyword,
			Page:    page,
		})
		fmt.Println(data)
		fmt.Println(err)
		if err != nil {
			ctx.JSON(resp.ServerErr(err, "查询失败"))
			return
		}
		ctx.JSON(resp.Success(data, "查询成功"))
	})

	// 搜索音乐结果详情
	g.GET("/search/:origin/:id", func(ctx *gin.Context) {
		origin := ctx.Param("origin")
		id := ctx.Param("id")

		data, err := service[origin].SearchDetail(id)
		if err != nil {
			ctx.JSON(resp.ServerErr(err, "查询失败"))
			return
		}
		ctx.JSON(resp.Success(data, "查询成功"))
	})

	// 音乐播放地址 返回音乐流
	g.GET("/music/file/:origin/:id", func(ctx *gin.Context) {
		origin := ctx.Param("origin")
		id := ctx.Param("id")

		proxy, req, err := service[origin].GetMusicFile(id)

		if err != nil {
			fmt.Printf("Err: %+v\n", err)
			ctx.JSON(resp.ServerErr(err, "获取歌曲文件失败"))
			return
		}

		proxy.ServeHTTP(ctx.Writer, req)
	})

	// 获取歌单广场歌单 后面传入源地址
	g.GET("/open-music-order", func(ctx *gin.Context) {
		originUrl := ctx.Query("origin")

		raw, err := http.Get(originUrl)

		defer raw.Body.Close()

		if err != nil || raw.StatusCode != 200 {
			ctx.JSON(resp.ServerErr(err, "请求失败"))
			return
		}

		if body, err := io.ReadAll(raw.Body); err != nil {
			ctx.JSON(resp.ServerErr(err, "数据读取失败"))
			return
		} else {
			result := []bb_type.MusicOrderItem{}
			if err := json.Unmarshal(body, &result); err != nil {
				ctx.JSON(resp.ServerErr(err, "json 序列化失败"))
				return
			}
			ctx.JSON(resp.Success(result, "请求成功"))
		}
	})

	// 图片代理服务
	g.GET("/img-proxy", func(ctx *gin.Context) {
		imgUrl := ctx.Query("url")
		u, err := url.Parse(imgUrl)
		if err != nil {
			ctx.JSON(resp.ServerErr(err, "图片地址错误请检查"))
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(u)

		req, err := http.NewRequest("GET", imgUrl, nil)
		if err != nil {
			ctx.JSON(resp.ServerErr(err, "图片请求出错"))
			return
		}
		proxy.ServeHTTP(ctx.Writer, req)
	})

	return bb_server.New(&http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}, log.Println)
}

// 自用应用log服务
type SvcLogger struct{}

func NewSvcLogger() *SvcLogger {
	return &SvcLogger{}
}

func (l *SvcLogger) Info(message ...string) {
	log.Println("BiliSvc Info | ", message)
}
func (l *SvcLogger) Warn(message ...string) {
	log.Println("BiliSvc Warn | ", message)
}
func (l *SvcLogger) Error(message ...string) {
	log.Println("BiliSvc Err | ", message)
}
