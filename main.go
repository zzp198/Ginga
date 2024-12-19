package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {

	GetOnceInfo()
	fmt.Println(getGPUInfo())
	fmt.Println(getIPs())
}

func StartServer() {

	g := gin.Default()

	g.GET("/state/self", func(c *gin.Context) {
		// 获取当前正在运行的 Goroutine 数量
		numGoroutines := runtime.NumGoroutine()
		// 获取操作系统的线程数
		numCPU := runtime.NumCPU()

		fmt.Printf("当前 Goroutine 数量: %d\n", numGoroutines)
		fmt.Printf("当前 CPU 核心数: %d\n", numCPU)

		p, err := process.NewProcess(int32(os.Getpid()))
		if err != nil {
			log.Fatal(err)
		}

		// 获取进程的内存信息
		memInfo, err := p.MemoryInfo()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Memory usage by current process: %v KB\n", memInfo.RSS/1024)

		// 获取当前 Goroutine 数量
		fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
	})

	g.Run(":80")
}

func GetOnceInfo() {

	fmt.Println(host.Info())
	//fmt.Println(host.SensorsTemperatures())
	//fmt.Println(host.Users())

}
func getGPUInfo() string {
	// 在Windows上获取GPU信息的方式依赖于使用的驱动程序
	cmd := exec.Command("wmic", "path", "win32_videocontroller", "get", "caption")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(output))
}

func getIPs() (string, string) {

	var internalIP string

	// 获取外网IP（使用公共API）
	// 这里我们简单地调用一个公共API来获取外网IP
	cmd := exec.Command("curl", "-s", "https://ipinfo.io/ip")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown", err.Error()
	}
	return internalIP, strings.TrimSpace(string(output))
}

func GetInfo() {

}
