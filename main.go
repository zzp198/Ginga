package main

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/mem"
	"io"
	"net/http"
)

func formatMemory(kb uint64) string {
	const (
		B = 1 << (10 * iota)
		KB
		MB
		GB
	)

	var value float64
	var unit string

	switch {
	case kb >= GB:
		value = float64(kb) / GB
		unit = "GB"
	case kb >= MB:
		value = float64(kb) / MB
		unit = "MB"
	case kb >= KB:
		value = float64(kb) / KB
		unit = "KB"
	case kb >= B:
		value = float64(kb) / KB
		unit = "B"
	default:
		value = float64(kb)
		unit = "B"
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

func main() {

	srv := gin.Default()

	apiGroup := srv.Group("api")
	apiGroup.GET("os_stat", func(c *gin.Context) {
		v, _ := mem.VirtualMemory()

		msg := fmt.Sprintf("Total: %s, Used:%s, UsedPercent:%f%%\n", formatMemory(v.Total), formatMemory(v.Used), v.UsedPercent)
		c.String(http.StatusOK, msg)
	})

	srv.GET("bilibili/:bv", func(c *gin.Context) {

		bv := c.Param("bv")

		view_url := "https://api.bilibili.com/x/web-interface/view?bvid=" + bv

		resp1, _ := http.Get(view_url)
		defer resp1.Body.Close()

		raw_data, _ := io.ReadAll(resp1.Body)

		raw_aid, _ := sonic.Get(raw_data, "data", "aid")
		raw_cid, _ := sonic.Get(raw_data, "data", "pages", 0, "cid")

		aid, _ := raw_aid.String()
		cid, _ := raw_cid.String()

		fmt.Println(aid)
		fmt.Println(cid)

		old_url := "https://api.bilibili.com/x/player/playurl?avid=" + aid + "&cid=" + cid + "&qn=80&fnval=1&platform=html5&high_quality=1"
		//new_url := "https://api.bilibili.com/x/player/wbi/playurl"

		fmt.Println(old_url)

		resp2, _ := http.Get(old_url)
		defer resp2.Body.Close()

		raw_data, _ = io.ReadAll(resp2.Body)

		url_node, _ := sonic.Get(raw_data, "data", "durl", 0, "url")

		url, _ := url_node.String()

		c.Redirect(http.StatusFound, url)
	})

	_ = srv.Run()
}
