package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"net/http"
	"zzp198/Ginga/util"
)

func main() {

	web := gin.Default()

	web.GET("/bili/:bv", func(c *gin.Context) {

		bv := c.Param("bv")

		view_request, err := util.HttpGet("https://api.bilibili.com/x/web-interface/view?bvid=" + bv)
		if err != nil {
			c.String(200, err.Error())
		}

		r1 := gjson.GetMany(view_request, "data.aid", "data.cid")
		aid := r1[0].String()
		cid := r1[1].String()

		fmt.Println(aid, cid)

		url := "https://api.bilibili.com/x/player/playurl?avid=" + aid + "&cid=" + cid + "&qn=80&fnval=1&platform=html5&high_quality=1"

		play_request, err := util.HttpGet(url)
		if err != nil {
			c.String(200, err.Error())
		}

		fmt.Println(play_request)
		real_url := gjson.Get(play_request, "data.durl.0.url").String()

		c.Redirect(http.StatusFound, real_url)
	})

	_ = web.Run(":80")
}
