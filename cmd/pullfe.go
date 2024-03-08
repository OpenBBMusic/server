package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Response struct {
	Assets []struct {
		Name string `json:"name"`
		Url  string `json:"browser_download_url"`
	}
}

func main() {
	url := getDownloadURL()
	fmt.Println("Download Url:", url)
	downloadFile(url)
}

func getDownloadURL() string {
	// 目标URL
	url := "https://api.github.com/repos/bb-music/web/releases/latest"

	// 使用http.Get获取数据
	resp, err := http.Get(url)
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
		if item.Name == "dist.tar.gz" {
			return item.Url
		}
	}
	return ""
}

func downloadFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// 创建并写入文件
	file, err := os.Create(filepath.Join("dist.tar.gz"))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	// 打印成功消息
	fmt.Println("File downloaded Success")
	return nil
}
