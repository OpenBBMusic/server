package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Response struct {
	Assets []struct {
		Name string `json:"name"`
		Url  string `json:"browser_download_url"`
	}
}

func main() {
	args := os.Args
	token := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--token" {
			token = args[i+1]
		}
	}
	url := getDownloadURL(token)
	fmt.Println("文件地址: ", url)
	unTarFile(url)
}

func getDownloadURL(token string) string {
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/bb-music/web/releases/latest", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	// 将响应内容解析为Go结构体
	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return ""
	}

	for _, item := range data.Assets {
		if item.Name == "web_dist.tar.gz" {
			return item.Url
		}
	}
	return ""
}

func unTarFile(url string) error {
	fmt.Println("下载前端资源")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("下载失败:", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println("下载完成，开始解压")
	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gr.Close()

	// 解压
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := filepath.Join(hdr.Name)
		fmt.Println("解压文件: ", filename)
		fileType := rune(hdr.Typeflag)
		if fileType != tar.TypeDir {
			file, err := createFile(filename)
			if err != nil {
				fmt.Println("解压失败", err)
				return err
			}
			io.Copy(file, tr)
		}
	}

	// 打印成功消息
	fmt.Println("解压完成")
	return nil
}

func createFile(name string) (*os.File, error) {
	name = strings.Replace(name, "client/web/", "", 1)
	dir := string([]rune(name)[0:strings.LastIndex(name, "/")])
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
