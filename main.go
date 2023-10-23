package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"time"
)

//go:embed frontend/* frontend/static/**/*
var StaticFS embed.FS

func main() {
	r := gin.Default()

	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(StaticFS, "frontend/*.html")))

	subfs, _ := fs.Sub(StaticFS, "frontend")
	r.StaticFS("/frontend", http.FS(subfs))

	r.GET("/Caddy/Start", func(c *gin.Context) {
		cmd := exec.Command("Caddy/Caddy", "start")

		msg, err := cmd.CombinedOutput()
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}

		c.String(http.StatusOK, string(msg))
	})

	r.GET("/Caddy/Stop", func(c *gin.Context) {
		cmd := exec.Command("Caddy/Caddy", "stop")
		cmd.WaitDelay = 10 * time.Second

		msg, err := cmd.CombinedOutput()
		if err != nil {
			c.String(http.StatusOK, err.Error())
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

		fiel, err := os.ReadFile("Caddy/Caddyfile")

		c.HTML(http.StatusOK, "Caddy.html", gin.H{
			"version":   string(version),
			"caddyfile": string(fiel),
		})
	})

	r.Run(":8080")
}
