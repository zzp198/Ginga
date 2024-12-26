package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {

	r := gin.Default()

	r.GET("/api/host", func(c *gin.Context) {
		hi, err := host.Info()
		if err != nil {
			c.JSON(200, err)
			return
		}

		pr, err := process.NewProcess(int32(os.Getpid()))
		if err != nil {
			c.JSON(200, err)
			return
		}

		// 获取进程的内存信息
		mi, err := pr.MemoryInfo()
		if err != nil {
			c.JSON(200, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"主机名":    hi.Hostname,
			"上线时间":   hi.Uptime,
			"操作系统":   hi.Platform,
			"系统版本":   hi.PlatformVersion,
			"内核数":    runtime.NumCPU(),
			"协程数":    runtime.NumGoroutine(),
			"物理内存占用": mi.RSS, // Resident Set Size
		})
	})

	r.GET("/api/host/state", func(c *gin.Context) {
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

	r.GET("/api/xray/versions", func(c *gin.Context) {
		rp, err := http.Get("https://api.github.com/repos/XTLS/xray-core/releases")
		if err != nil {
			c.JSON(http.StatusOK, err)
			return
		}

		b, err := io.ReadAll(rp.Body)
		_ = rp.Body.Close()

		var r []string
		for _, name := range gjson.GetBytes(b, "#.tag_name").Array() {
			r = append(r, name.String())
		}
		c.JSON(http.StatusOK, r)
	})

	r.GET("/api/xray/changeversion", func(c *gin.Context) {

	})
	_ = r.Run(":8880")
}
