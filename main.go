package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed frontend/* frontend/static/**/*
var StaticFS embed.FS

func main() {

	ip := flag.String("ip", "0.0.0.0:8080", "")

	r := gin.Default()

	r.LoadHTMLFiles("../frontend/*.html")

	frontend, _ := fs.Sub(StaticFS, "frontend")
	r.StaticFS("/frontend", http.FS(frontend))

	r.GET("/api", func(c *gin.Context) {

		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Failed to read message: ", err)
				break
			}

			fmt.Println("Received message: %s\n", string(msg))

			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("Failed to write message: ", err)
				break
			}
		}

	})

	_ = r.Run(*ip)
}
