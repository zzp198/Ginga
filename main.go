package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	thirednet "github.com/shirou/gopsutil/net"
	"net/http"
	"time"
)

func main() {

}

func main1() {

	//frontend.Server("")

	Web := gin.Default()

	Web.GET("/events", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		for {
			fmt.Fprintf(c.Writer, "data: %s\n\n", time.Now().Format(time.Stamp))
			c.Writer.(http.Flusher).Flush()
			time.Sleep(1 * time.Second)
		}
	})

	Web.GET("/state", func(c *gin.Context) {
		State()
		c.String(200, "")
	})

	Web.GET("/load", func(c *gin.Context) {

	})

	Web.Run(":8080")

	//MailServer()

}

func State() {

	// CPU 信息
	cpuInfo, _ := cpu.Info()
	fmt.Println("CPU Info:", cpuInfo)

	precents, _ := cpu.Percent(0, false)
	fmt.Println("Precents:", precents)

	uptime, _ := host.Uptime()
	fmt.Println("Uptime:", uptime)

	// 内存信息
	memInfo, _ := mem.VirtualMemory()
	fmt.Println("Memory Info:", memInfo)

	swapInfo, _ := mem.SwapMemory()
	fmt.Println("Swap Memory Info:", swapInfo)

	// 磁盘信息
	diskInfo, _ := disk.Usage("/")
	fmt.Println("Disk Info:", diskInfo)

	// 网络信息
	netInfo, _ := thirednet.Interfaces()
	fmt.Println("Network Info:", netInfo)

	avgState, _ := load.Avg()
	fmt.Println("Load Average State:", avgState)

	connState, _ := thirednet.Connections("all")
	fmt.Println("connState:", connState, len(connState))

	ioState, _ := thirednet.IOCounters(false)
	fmt.Println("ioState:", ioState)
}
