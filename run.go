package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"zzp198/Ginga/util"
)

type ServerInfo struct {
	gorm.Model
	Ip       string
	Username string
	Password string
	Key      string
	ViewTime int
}

func main() {

	db := SqliteConn()

	err := db.AutoMigrate(&ServerInfo{})
	if err != nil {
		panic(err)
	}

	web := gin.Default()

	web.GET("/bili/:bv", func(c *gin.Context) {

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

	web.GET("/os_stat/", func(c *gin.Context) {
		v, _ := mem.VirtualMemory()

		msg := fmt.Sprintf("Total: %s, Used:%s, UsedPercent:%f%%\n", util.FormatByte(v.Total), util.FormatByte(v.Used), v.UsedPercent)
		c.String(http.StatusOK, msg)
	})

	web.GET("/server/", func(c *gin.Context) {
		var results []ServerInfo

		db.Find(&results)

		var msg string
		for _, product := range results {
			msg += fmt.Sprintf("IP: %s, USER: %s, PASS: %s\n", product.Ip, product.Username, product.Password)
		}

		c.String(200, msg)
	})

	_ = web.Run(":80")
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
