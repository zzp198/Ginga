package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	g := gin.New()

	g.GET("/", func(c *gin.Context) {

		cmd := exec.Command("/bin/bash caddy version")
		cmd.Run()

		c.HTML(200, "showcaddy.html",
			gin.H{
				"exp": "1",
			})
	})

	g.GET("/downloadcaddy", func(c *gin.Context) {
		version, exist := c.GetQuery("version")
		if !exist {
			c.AbortWithStatus(404)
			return
		}

		filename, err := DownloadCaddy(version)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		// 打开 tar.gz 文件
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Failed to open file: %v\n", err)
			return
		}
		defer file.Close()

		os.Mkdir("caddy", 0755)

		// 提取到指定的目录
		err = extractTarGz(file, "caddy")
		if err != nil {
			fmt.Printf("Failed to extract tar.gz: %v\n", err)
			return
		}

		c.String(200, "success")
	})

	g.GET("/caddyversion", func(c *gin.Context) {

		beta := c.DefaultQuery("beta", "0")

		versions, err := GetCaddyVersion(beta)
		if err != nil {
			c.AbortWithError(500, err)
		}

		c.JSON(200, versions)
	})

	g.GET("/changeversion", func(c *gin.Context) {

	})

	g.GET("/start", func(c *gin.Context) {

	})

	g.GET("/stop", func(c *gin.Context) {

	})

	g.GET("/restart", func(c *gin.Context) {

	})

	g.Run(":80")
}

func GetCaddyVersion(beta string) ([]string, error) {

	resp, err := newHTTPClientWithProxy().Get("https://api.github.com/repos/caddyserver/caddy/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type TagName struct {
		TagName string `json:"name"`
	}

	var tag []TagName

	err = json.Unmarshal(body, &tag)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(tag))
	for _, v := range tag {
		if beta != "1" && strings.Contains(v.TagName, "beta") {
			continue
		}
		versions = append(versions, v.TagName)
	}

	return versions, nil

}

func DownloadCaddy(version string) (string, error) {

	filename := fmt.Sprintf("caddy_%s_linux_amd64.tar.gz", strings.TrimPrefix(version, "v"))
	url := fmt.Sprintf("https://github.com/caddyserver/caddy/releases/download/%s/%s", version, filename)

	resp, err := newHTTPClientWithProxy().Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	os.Remove(filename)

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// 解压 tar.gz 文件
func extractTarGz(gzipStream io.Reader, targetDir string) error {
	// 创建 gzip.Reader
	gzipReader, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// 创建 tar.Reader
	tarReader := tar.NewReader(gzipReader)

	// 循环读取 tar 包中的文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// 读取完成
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// 获取文件路径
		filePath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(filePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// 创建文件
			outFile, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer outFile.Close()

			// 将内容写入文件
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to write file content: %w", err)
			}
		default:
			// 跳过其他类型
			fmt.Printf("Skipping unknown type: %c in %s\n", header.Typeflag, filePath)
		}
	}
	return nil
}

func newHTTPClientWithProxy() *http.Client {
	// 获取系统环境变量中的代理地址
	//proxyURL, err := http.ProxyFromEnvironment(&http.Request{
	//	URL: &url.URL{},
	//})
	//if err != nil {
	//	fmt.Printf("Failed to get proxy from environment: %v\n", err)
	//	return &http.Client{}
	//}

	//fmt.Println(proxyURL)

	url, _ := url.Parse("http://127.0.0.1:10809") // v2rayN的自动配置系统代理没法获取

	// 创建带有代理设置的 Transport
	transport := &http.Transport{
		//Proxy: http.ProxyURL(proxyURL),
		Proxy: http.ProxyURL(url),
	}

	// 创建自定义 HTTP 客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 设置超时时间
	}
	return client
}
