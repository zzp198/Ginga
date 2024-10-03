package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	r := gin.Default()

	// 定义要浏览的根目录，比如 / 或 /home/user 等
	rootPath := "/"

	// 文件浏览路由
	r.GET("/files/*filepath", func(c *gin.Context) {
		path := c.Param("filepath")
		if path == "/" {
			path = ""
		}
		// 组合访问的路径与根目录路径
		fullPath := filepath.Join(rootPath, path)

		// 检查路径是否存在
		fileInfo, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "Path not found")
			return
		}

		// 如果是目录，列出目录内容
		if fileInfo.IsDir() {
			files, err := os.ReadDir(fullPath)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error reading directory")
				return
			}

			var fileList []string
			for _, file := range files {
				if file.IsDir() {
					fileList = append(fileList, file.Name()+"/") // 用 '/' 标记文件夹
				} else {
					fileList = append(fileList, file.Name())
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"current_dir": path,
				"files":       fileList,
			})
		} else {
			// 如果是文件，直接下载文件
			c.File(fullPath)
		}
	})

	// 启动服务器
	r.Run(":8080")
}
