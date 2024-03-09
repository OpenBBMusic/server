package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/bb-music/desktop/app_bili"
	"github.com/bb-music/desktop/pkg/bb_server"
	"github.com/bb-music/desktop/pkg/bb_type"
	"github.com/bb-music/server/internal/logger"
	"github.com/bb-music/server/internal/resp"
	"github.com/gin-gonic/gin"
)

type OriginService interface {
	GetConfig() any
	InitConfig(force bool) error
	Search(params bb_type.SearchParams) (*bb_type.SearchResponse, error)
	SearchDetail(id string) (*bb_type.SearchItem, error)
	GetMusicFile(id string) (*httputil.ReverseProxy, *http.Request, error)
	DownloadMusic(params bb_type.DownloadMusicParams) (string, error)
}

func NewServer(r *gin.Engine, port int, cacheDir string) *bb_server.Server {
	log.Println("======= API 服务启动 =======")
	log.Println("端口:", port)
	log.Println("缓存目录:", cacheDir)
	log.Println("注册音乐源服务")
	// 日志输出
	svcLogger := logger.NewSvcLogger()

	bili := app_bili.New(cacheDir, svcLogger)
	// 注册源服务
	service := map[bb_type.OriginType]OriginService{
		bb_type.BiliOriginName: bili,
	}
	// 获取源的配置信息
	bili.InitConfig(false)

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

	return bb_server.New(&http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}, log.Println)
}
