package main

import (
	"example/common"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var jwt string

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("resource/*")

	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			c.Next()
		} else {
			n, err := c.Cookie("jwt")
			if err == nil && n == jwt {
				c.Next()
			} else {
				c.Redirect(http.StatusTemporaryRedirect, "/")
				c.Abort()
			}
		}
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/", func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "admin" && password == "123456" {

			jwt = time.Now().String()

			c.SetCookie("jwt", jwt, 86400, "/", "", false, false)

			// 返回 JSON 数据，表示登录成功
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "登录成功",
			})
			return
		}

		// 登录失败，返回错误信息
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "用户名或密码错误",
		})
	})

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	r.GET("/api/info/static", func(c *gin.Context) {
		hostInfo, err := host.Info()

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  "-1",
				"error": err.Error(),
			})
			return
		}

		p, err := process.NewProcess(int32(os.Getpid()))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  "-1",
				"error": err.Error(),
			})
			return
		}

		// 获取进程的内存信息
		memInfo, err := p.MemoryInfo()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  "-1",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "ok",

			"主机名":    hostInfo.Hostname,
			"上线时间":   hostInfo.Uptime,
			"操作系统":   hostInfo.Platform,
			"系统版本":   hostInfo.PlatformVersion,
			"内核数":    runtime.NumCPU(),
			"协程数":    runtime.NumGoroutine(),
			"物理内存占用": common.FormatByte(memInfo.RSS), // Resident Set Size
		})
	})

	r.GET("/api/info/dynamic", func(c *gin.Context) {
		cmd := exec.Command("curl", "-s", "https://ipinfo.io/ip")
		output, err := cmd.Output()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  "-1",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": "1",
			"ip":   strings.TrimSpace(string(output)),
		})
	})

	r.Run(":80")
}
