package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
)

//go:embed frontend/*
var Frontend embed.FS

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("frontend/**/*")
	r.StaticFS("/frontend", http.FS(Frontend))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/caddy/status", func(c *gin.Context) {
		cmd := exec.Command("Caddy/Caddy", "-v")
		msg, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err,
			})
			return
		}
		c.String(http.StatusOK, string(msg))
	})

	r.GET("/Caddy", func(c *gin.Context) {
		cmd := exec.Command("Caddy/Caddy", "-v")
		version, err := cmd.CombinedOutput()
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}

		c.HTML(http.StatusOK, "Caddy.html", gin.H{
			"version": version,
		})
	})

	r.Run(":8080")
}
