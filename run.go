package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
	"runtime"
	"time"
	"zzp198/Ginga/util"
)

type ServerInfo struct {
	ID       uint `gorm:"primarykey"`
	Ip       string
	Username string
	Password string
	Key      string
	ViewTime int
}

type MemoryRecord struct {
	UsedPercent float64
	RecodeTime  int64
}

var ip *string

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//if token == "" {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		//	c.Abort()
		//	return
		//}

		// TODO Check Token

		c.Next()
	}
}

func main() {
	ip = flag.String("ip", ":8080", "ip address")
	flag.Parse()

	db := SqliteConn()

	err := db.AutoMigrate(&ServerInfo{}, &MemoryRecord{})
	if err != nil {
		panic(err)
	}

	c := cron.New()
	_, _ = c.AddFunc("@every 1m", func() {
		vm, err := mem.VirtualMemory()
		if err != nil {
		}

		db.Create(&MemoryRecord{
			UsedPercent: vm.UsedPercent,
			RecodeTime:  time.Now().Unix(),
		})
	})
	c.Start()

	r := gin.Default()
	r.Use(AuthMiddleware())

	r.GET("/", func(c *gin.Context) {

	})

	r.GET("/bili/:bv", func(c *gin.Context) {

		bv := c.Param("bv")

		view_request, err := util.HttpGet("https://api.bilibili.com/x/web-interface/view?bvid=" + bv)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		r1 := gjson.GetMany(view_request, "data.aid", "data.cid")
		aid := r1[0].String()
		cid := r1[1].String()

		fmt.Println(aid, cid)

		url := "https://api.bilibili.com/x/player/playurl?avid=" + aid + "&cid=" + cid + "&qn=80&fnval=1&platform=html5&high_quality=1"

		play_request, err := util.HttpGet(url)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		fmt.Println(play_request)
		real_url := gjson.Get(play_request, "data.durl.0.url").String()

		c.Redirect(http.StatusFound, real_url)
	})

	r.GET("/api/os_stat/", func(c *gin.Context) {
		v, _ := mem.VirtualMemory()

		c.JSON(http.StatusOK, gin.H{
			"Total":       util.FormatBit(v.Total),
			"Available":   util.FormatBit(v.Available),
			"UsedPercent": v.UsedPercent,
		})
	})

	r.GET("/server/", func(c *gin.Context) {
		var results []ServerInfo

		db.Find(&results)

		var msg string
		for _, product := range results {
			msg += fmt.Sprintf("IP: %s, USER: %s, PASS: %s\n", product.Ip, product.Username, product.Password)
		}

		c.String(200, msg)
	})

	r.GET("/api/self", func(c *gin.Context) {

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		runtime.NumGoroutine()

		pid := os.Getpid()
		p, err := process.NewProcess(int32(pid))
		if err != nil {
			c.Abort()
		}

		meminfo, err := p.MemoryInfo()
		if err != nil {
			c.Abort()
		}

		c.JSON(200, gin.H{
			"RSS":          util.FormatBit(meminfo.RSS * 8), // 常驻内存
			"VMS":          util.FormatBit(meminfo.VMS * 8), // 虚拟内存
			"Alloc":        util.FormatBit(m.Alloc * 8),
			"TotalAlloc":   util.FormatBit(m.TotalAlloc * 8),
			"Sys":          util.FormatBit(m.Sys * 8),
			"NumGC":        m.NumGC,
			"NumGoroutine": runtime.NumGoroutine(),
		})
	})

	r.GET("/api/cpu", func(c *gin.Context) {
		var results []MemoryRecord
		db.Order("recode_time DESC").Limit(10).Find(&results)
		c.JSON(200, results)
	})

	r.GET("api/net", func(c *gin.Context) {

		now1 := time.Now()
		netstat1, err := net.IOCounters(false)
		if err != nil || len(netstat1) < 1 {
			c.Abort()
		}

		time.Sleep(1 * time.Second)

		now2 := time.Now()
		netstat2, err := net.IOCounters(false)
		if err != nil || len(netstat2) < 1 {
			c.Abort()
		}

		seconds := now2.Sub(now1).Seconds()
		up := float64(netstat2[0].BytesSent-netstat1[0].BytesSent) / seconds
		down := float64(netstat2[0].BytesRecv-netstat1[0].BytesRecv) / seconds

		c.JSON(200, gin.H{
			"up":   util.FormatBit(uint64(up * 8)),
			"down": util.FormatBit(uint64(down * 8)),
		})
	})

	r.GET("api/xrayv", func(c *gin.Context) {

		rbody, err := util.HttpGet("https://api.github.com/repos/XTLS/Xray-core/releases")
		if err != nil {
			c.Abort()
		}

		var tags []struct {
			Tag string `json:"tag_name"`
		}

		json.Unmarshal([]byte(rbody), &tags)

		c.JSON(200, tags)
	})

	r.GET("api/load", func(c *gin.Context) {

		stst, err := load.Avg()
		if err != nil {
			c.Abort()
		}

		c.JSON(200, gin.H{
			"avg1":  stst.Load1,
			"avg5":  stst.Load5,
			"avg15": stst.Load15,
		})
	})

	_ = r.Run(*ip)
}

func SqliteConn() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("ginga.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func MysqlConn(host, port, user, pass, dbname string) *gorm.DB {
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(fmt.Sprintf(dsn, user, pass, host, port, dbname)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db

}
