package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"net/http"
)

func main() {

	State()

	fmt.Println("1")

	Web := gin.Default()

	Web.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	Web.Run(":80")
}

func State() {

	// CPU 信息
	cpuInfo, _ := cpu.Info()
	fmt.Println("CPU Info:", cpuInfo)

	// 内存信息
	memInfo, _ := mem.VirtualMemory()
	fmt.Println("Memory Info:", memInfo)

	// 磁盘信息
	diskInfo, _ := disk.Usage("/")
	fmt.Println("Disk Info:", diskInfo)

	// 网络信息
	netInfo, _ := net.Interfaces()
	fmt.Println("Network Info:", netInfo)

}
