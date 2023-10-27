package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"net/http"
	"time"
)

func State(c *gin.Context) {
	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)
	fmt.Println()

	vm, _ := mem.VirtualMemory()

	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)

	//totalPercent, _ := cpu.Percent(1*time.Second, false) // 有时延
	//perPercents, _ := cpu.Percent(1*time.Second, true)

	//cpuinfos, _ := cpu.Info()

	version, _ := host.KernelVersion()

	platform, family, version, _ := host.PlatformInformation()

	diskinfos, _ := disk.Partitions(false)

	diskusage, _ := disk.Usage("/")

	swapMemory, _ := mem.SwapMemory()

	c.JSON(http.StatusOK, gin.H{
		"物理核数": physicalCnt,
		"逻辑核数": logicalCnt,

		"VirtualMemoryAvailable": vm,
		"BootTime":               t.Local().Format("2006-01-02 15:04:05"),
		"platform":               platform,
		"family":                 family,
		"version":                version,
		//"totalPercent":           totalPercent,
		//"perPercents":            perPercents,
		//"cpuinfos:": infos,// flags很多
		"diskinfos":  diskinfos,
		"diskusage":  diskusage,
		"swapMemory": swapMemory,
	})
}
